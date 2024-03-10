# credit to https://github.com/tryolabs/norfair/blob/master/demos/yolov7/src/demo.py for the basis of this code

import datetime
import traceback
import os
import time
from concurrent.futures import ThreadPoolExecutor, as_completed
from threading import RLock
from typing import Callable, List, Union, Optional, NamedTuple, Dict, Hashable, Tuple

import numpy as np
import torch
import torch.backends.cuda
from cv2 import CAP_PROP_POS_MSEC, CAP_PROP_FPS  # noqa
from norfair import Detection, Tracker, draw_points, draw_boxes
from norfair.tracker import TrackedObject
from orjson import dumps

from .video import Video
from .helpers import DetectionContext

DISTANCE_THRESHOLD_BBOX: float = 7.0
DISTANCE_THRESHOLD_CENTROID: int = 30
MAX_DISTANCE: int = 10000
INITIALIZATION_DELAY: int = 2

CONF_THRESHOLD: float = 0.44
IOU_THRESHOLD: float = 0.1
IMAGE_SIZE: int = 640
STRIDE_FRAMES: int = 10


class RawDetectedObject(NamedTuple):
    tracking_mode: str
    object_id: int
    class_id: Union[int, Optional[Hashable]]
    class_name: str
    timedeltas: List[datetime.timedelta]


class DetectedObject(NamedTuple):
    object_id: int
    class_id: Union[int, Optional[Hashable]]
    class_name: str
    start_timedelta: datetime.timedelta
    end_timedelta: datetime.timedelta
    duration_timedelta: datetime.timedelta


class ProcessedVideo(NamedTuple):
    input_path: str
    output_path: str
    stride_frames: int
    frames_per_second: int
    duration: datetime.timedelta
    total_frames: int
    handled_frames: int
    detected_objects: List[DetectedObject]

    def to_json(self) -> bytes:
        return dumps(
            {
                "input_path": self.input_path,
                "output_path": self.output_path,
                "stride_frames": self.stride_frames,
                "frames_per_second": self.frames_per_second,
                "duration": self.duration.total_seconds(),
                "total_frames": self.total_frames,
                "handled_frames": self.handled_frames,
                "detected_objects": [
                    {
                        "object_id": detected_object.object_id,
                        "class_id": detected_object.class_id,
                        "class_name": detected_object.class_name,
                        "start_timedelta": detected_object.start_timedelta.total_seconds(),
                        "end_timedelta": detected_object.end_timedelta.total_seconds(),
                        "duration_timedelta": detected_object.duration_timedelta.total_seconds(),
                    }
                    for detected_object in self.detected_objects
                ],
            }
        )


def yolo_detections_to_norfair_detections(
    yolo_detections: torch.tensor,
) -> Tuple[List[Detection], List[Detection]]:
    centroid_norfair_detections: List[Detection] = []
    bbox_norfair_detections: List[Detection] = []

    detections_as_xywh = yolo_detections.xywh[0]
    for detection_as_xywh in detections_as_xywh:
        centroid = np.array(
            [
                detection_as_xywh[0].item(),
                detection_as_xywh[1].item(),
            ]
        )

        scores = np.array(
            [
                detection_as_xywh[4].item(),
            ]
        )

        detection = Detection(
            points=centroid,
            scores=scores,
            label=int(detection_as_xywh[-1].item()),
        )

        centroid_norfair_detections.append(detection)

    detections_as_xyxy = yolo_detections.xyxy[0]
    for detection_as_xyxy in detections_as_xyxy:
        bbox = np.array(
            [
                [
                    detection_as_xyxy[0].item(),
                    detection_as_xyxy[1].item(),
                ],
                [
                    detection_as_xyxy[2].item(),
                    detection_as_xyxy[3].item(),
                ],
            ]
        )

        scores = np.array(
            [
                detection_as_xyxy[4].item(),
                detection_as_xyxy[4].item(),
            ]
        )

        detection = Detection(
            points=bbox,
            scores=scores,
            label=int(detection_as_xyxy[-1].item()),
        )

        bbox_norfair_detections.append(detection)

    return centroid_norfair_detections, bbox_norfair_detections


class ObjectTracker(object):
    def __init__(
        self,
        model_name: str,
        device: str,
        tracking_mode: str,
        conf_threshold: float,
        iou_threshold: float,
        img_size: int,
        stride_frames: int,
        handle_detections: Callable[[DetectionContext], None],
    ):
        model_path = os.path.join(
            os.path.abspath(os.path.split(__file__)[0]), "models", model_name
        )
        if not os.path.exists(model_path):
            raise ValueError(f"model_path={repr(model_path)} does not exist.")

        device = device or "cuda"

        self._executor = ThreadPoolExecutor()
        self._futures = []

        self._model = torch.hub.load("WongKinYiu/yolov7", "custom", model_path)
        self._model.to(torch.device(device))

        self._model.conf = conf_threshold
        self._model.iou = iou_threshold

        self._tracking_mode = tracking_mode
        self._img_size = img_size
        self._stride_frames = stride_frames
        self._handle_detections = handle_detections

        self._name_by_class_id = {i: x for i, x in enumerate(self._model.names)}
        self._lock = RLock()
        self._detected_object_by_object_id: Dict[int, RawDetectedObject] = {}

    def __del__(self):
        print("tearing down {}".format(self._executor))
        try:
            self._executor.shutdown(wait=True)
        except Exception:
            pass

        print("tearing down {}".format(self._model))
        try:
            del self._model
        except Exception:
            pass

        print("tearing down {}".format(self))
        try:
            super().__del__()
        except Exception:
            pass
        try:
            del self
        except Exception:
            pass

    def _draw_tracking_context_on_frame(
        self,
        frame: np.ndarray,
        tracked_objects: List[TrackedObject],
        video: Video,
    ):
        if self._tracking_mode == "centroid":
            draw_points(
                frame=frame,
                drawables=tracked_objects,
                color="by_id",
                thickness=1,
                draw_labels=True,
                text_size=1,
                text_thickness=1,
                draw_ids=True,
                draw_scores=True,
            )

        elif self._tracking_mode == "bbox":
            draw_boxes(
                frame,
                tracked_objects,
                color="by_id",
                thickness=1,
                draw_labels=True,
                text_size=1,
                text_thickness=1,
                draw_ids=True,
                draw_scores=True,
            )

        video.write(frame)

    def _handle_yolo_detections(
        self,
        tracker: Tracker,
        yolo_detections: torch.tensor,
        frame: np.ndarray,
        video: Video,
        video_timedelta: datetime.timedelta,
    ) -> Tuple[float, float]:
        centroid_detections, bbox_detections = yolo_detections_to_norfair_detections(
            yolo_detections=yolo_detections
        )

        self._handle_detections(
            detection_context=DetectionContext(
                centroid_detections=centroid_detections,
                bbox_detections=bbox_detections,
                timedelta=video_timedelta,
            ),
            name_by_class_id=self._name_by_class_id,
        )

        tracked_objects: List[TrackedObject] = tracker.update(
            detections=(
                centroid_detections
                if self._tracking_mode == "centroid"
                else bbox_detections
            ),
        )

        self._draw_tracking_context_on_frame(
            frame=frame,
            tracked_objects=tracked_objects,
            video=video,
        )

        for tracked_object in tracked_objects:
            with self._lock:
                detected_object = self._detected_object_by_object_id.setdefault(
                    tracked_object.global_id,
                    RawDetectedObject(
                        tracking_mode=self._tracking_mode,
                        object_id=tracked_object.global_id,
                        class_id=tracked_object.last_detection.label,
                        class_name=self._name_by_class_id.get(
                            tracked_object.last_detection.label
                        )
                        or "__unknown__",
                        timedeltas=[],
                    ),
                )

                detected_object.timedeltas.append(video_timedelta)

    def _handle_file(
        self,
        input_path: str,
        output_path: Optional[str] = None,
    ) -> ProcessedVideo:
        before = time.time()

        if not output_path:
            folder_path, file_name = os.path.split(input_path)
            file_base_name, file_extension = os.path.splitext(file_name)
            output_path = os.path.join(
                folder_path,
                f"{file_base_name}_out.{file_extension.lstrip('.')}",
            )

        if not os.path.exists(input_path) or not os.path.isfile(input_path):
            raise ValueError(
                f"input_path={repr(input_path)} does not exist or is not a file"
            )

        print("setting up video and tracker for path={}".format(repr(input_path)))

        video = Video(
            input_path=os.path.abspath(input_path),
            output_path=os.path.abspath(output_path),
        )

        frames_per_second: int = video.video_capture.get(CAP_PROP_FPS)

        # TOOD: dear future self- attempts to reduce this below 1 fps have been in vain
        # video.output_fps = float(frames_per_second) / float(self._stride_frames)

        distance_function = "iou" if self._tracking_mode == "bbox" else "euclidean"
        distance_threshold = (
            DISTANCE_THRESHOLD_BBOX
            if self._tracking_mode == "bbox"
            else DISTANCE_THRESHOLD_CENTROID
        )
        distance_threshold *= self._stride_frames

        tracker = Tracker(
            distance_function=distance_function,
            distance_threshold=distance_threshold,
            initialization_delay=INITIALIZATION_DELAY,
        )

        self._detected_object_by_object_id = {}

        video_timedelta = datetime.timedelta(seconds=0)
        total_frames = 0
        handled_frames = 0

        print("iterating frames...")

        for i, frame in enumerate(video):
            total_frames += 1

            video_timedelta = datetime.timedelta(
                milliseconds=video.video_capture.get(CAP_PROP_POS_MSEC)
            )

            if i % self._stride_frames != 0:
                continue

            handled_frames += 1

            yolo_detections: torch.tensor = self._model(frame, size=self._img_size)

            self._futures.append(
                self._executor.submit(
                    self._handle_yolo_detections,
                    tracker=tracker,
                    yolo_detections=yolo_detections,
                    frame=frame,
                    video=video,
                    video_timedelta=video_timedelta,
                )
            )

        print("waiting for drawing futures...")

        for future in as_completed(self._futures):
            try:
                _ = future.result()
            except Exception:
                traceback.print_exc()
                raise

        with self._lock:
            raw_detected_objects: List[RawDetectedObject] = sorted(
                list(self._detected_object_by_object_id.values()),
                key=lambda x: x.object_id,
            )

            detected_objects = [
                DetectedObject(
                    object_id=x.object_id,
                    class_id=x.class_id,
                    class_name=x.class_name,
                    start_timedelta=min(x.timedeltas),
                    end_timedelta=max(x.timedeltas),
                    duration_timedelta=max(x.timedeltas) - min(x.timedeltas),
                )
                for x in raw_detected_objects
            ]

            processed_video = ProcessedVideo(
                input_path=input_path,
                output_path=output_path,
                stride_frames=self._stride_frames,
                frames_per_second=frames_per_second,
                duration=video_timedelta,
                total_frames=total_frames,
                handled_frames=handled_frames,
                detected_objects=detected_objects,
            )
            self._detected_object_by_object_id = {}
            self._futures = []

        after = time.time()

        print(
            "duration={} for processed_video={}".format(
                after - before, processed_video.to_json()
            )
        )

        return processed_video

    def __call__(
        self,
        input_path: str,
        output_path: Optional[str] = None,
    ) -> ProcessedVideo:
        return self._handle_file(
            input_path=input_path,
            output_path=output_path,
        )
