import datetime
import os
import time
import signal

from typing import Callable, cast, Any, List
from threading import Event
from functools import partial

import cv2
import numpy as np
import torch

from .types import class_name_by_class_id, Detection, DetectedFrame


IMAGE_SIZE: float = 640.0
STRIDE_FRAMES: float = 4.0

_STOP = Event()


def signal_handler(sig, frame):
    _STOP.set()


def get_model(raw_model_path: str) -> Any:
    model_path = os.path.abspath(raw_model_path)

    if not os.path.exists(raw_model_path):
        raise ValueError(f"{repr(raw_model_path)} does not exist")

    if not os.path.isfile(raw_model_path):
        raise ValueError(f"{repr(raw_model_path)} exists but it is not a file")

    device = None
    if torch.cuda.is_available():
        device = torch.device("cuda:0")
        torch.cuda.synchronize()
    elif torch.backends.mps.is_available():
        device = torch.device("mps")
    else:
        device = torch.device("cpu")
    print(f"device: {device}")

    model = torch.hub.load("WongKinYiu/yolov7", "custom", model_path)
    model.to(device)
    print(f"model: {model}")

    return model


def handle_frame(
    model: Any,
    frame: np.array,
    current_frame: float,
    fps: float,
) -> DetectedFrame:
    timedelta = datetime.timedelta(seconds=current_frame / fps)

    yolo_detections: torch.tensor = model(frame, size=IMAGE_SIZE)

    detections: List[Detection] = []

    for xyxys in yolo_detections.xyxy:
        for xyxy in xyxys:
            detections.append(Detection.from_xyxy(xyxy, current_frame, timedelta))

    return DetectedFrame(
        frame=current_frame,
        timedelta=timedelta,
        detections=detections,
    )


def handle_video(
    raw_file_path: str,
    on_frame: Callable[[], None],
) -> None:
    try:
        file_path = os.path.abspath(raw_file_path)

        if not os.path.exists(raw_file_path):
            raise ValueError(f"{repr(raw_file_path)} does not exist")

        if not os.path.isfile(raw_file_path):
            raise ValueError(f"{repr(raw_file_path)} exists but it is not a file")

        video_capture = cv2.VideoCapture(file_path)
        print(f"{video_capture}")

        total_frames = float(video_capture.get(cv2.CAP_PROP_FRAME_COUNT))
        print(f"total_frames: {total_frames}")

        fps = float(video_capture.get(cv2.CAP_PROP_FPS))
        print(f"fps: {fps}")

        width = float(video_capture.get(cv2.CAP_PROP_FRAME_WIDTH))
        print(f"width: {width}")

        height = float(video_capture.get(cv2.CAP_PROP_FRAME_HEIGHT))
        print(f"height: {height}")

        start = time.time()
        current_frame = 0.0
        fpses = []
        detected_objects = 0

        while not _STOP.is_set():
            current_frame += 1

            ret, frame = video_capture.read()

            ret: bool = cast(bool, ret)
            frame: np.array = cast(np.array, frame)

            if ret is False or frame is None:
                break

            if current_frame % STRIDE_FRAMES != 0:
                continue

            detected_frame: DetectedFrame = on_frame(
                frame=frame,
                current_frame=current_frame,
                fps=fps,
            )

            detected_objects += len(detected_frame.detections)

            now = time.time()

            fpses.append(float(current_frame) / (now - start))

        handled_frames = len(fpses)

        min_fps = 0
        avg_fps = 0
        max_fps = 0

        if handled_frames > 0:
            min_fps = min(fpses)
            avg_fps = sum(fpses) / handled_frames
            max_fps = max(fpses)

        print(f"handled_frames: {handled_frames}")
        print(f"min_fps: {min_fps}")
        print(f"avg_fps: {avg_fps}")
        print(f"max_fps: {max_fps}")
        print(f"detected_objects: {detected_objects}")
    finally:
        cv2.destroyAllWindows()


def setup():
    signal.signal(signal.SIGINT, signal_handler)

    raw_model_path = os.getenv("MODEL_PATH") or ""
    if not raw_model_path:
        raise ValueError("MODEL_PATH env var empty or unset")

    model = get_model(
        raw_model_path=raw_model_path,
    )


def run_once():
    signal.signal(signal.SIGINT, signal_handler)

    raw_model_path = os.getenv("MODEL_PATH") or ""
    if not raw_model_path:
        raise ValueError("MODEL_PATH env var empty or unset")

    raw_file_path = os.getenv("FILE_PATH") or ""
    if not raw_file_path:
        raise ValueError("FILE_PATH env var empty or unset")

    model = get_model(
        raw_model_path=raw_model_path,
    )

    handle_video(
        raw_file_path,
        partial(
            handle_frame,
            model=model,
        ),
    )


def run():
    signal.signal(signal.SIGINT, signal_handler)

    print("Hello, world")
