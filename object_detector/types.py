import datetime

from typing import NamedTuple, Any, List


class_name_by_class_id = dict(
    enumerate(
        [
            "person",
            "bicycle",
            "car",
            "motorcycle",
            "airplane",
            "bus",
            "train",
            "truck",
            "boat",
            "traffic light",
            "fire hydrant",
            "stop sign",
            "parking meter",
            "bench",
            "bird",
            "cat",
            "dog",
            "horse",
            "sheep",
            "cow",
            "elephant",
            "bear",
            "zebra",
            "giraffe",
            "backpack",
            "umbrella",
            "handbag",
            "tie",
            "suitcase",
            "frisbee",
            "skis",
            "snowboard",
            "sports ball",
            "kite",
            "baseball bat",
            "baseball glove",
            "skateboard",
            "surfboard",
            "tennis racket",
            "bottle",
            "wine glass",
            "cup",
            "fork",
            "knife",
            "spoon",
            "bowl",
            "banana",
            "apple",
            "sandwich",
            "orange",
            "broccoli",
            "carrot",
            "hot dog",
            "pizza",
            "donut",
            "cake",
            "chair",
            "couch",
            "potted plant",
            "bed",
            "dining table",
            "toilet",
            "tv",
            "laptop",
            "mouse",
            "remote",
            "keyboard",
            "cell phone",
            "microwave",
            "oven",
            "toaster",
            "sink",
            "refrigerator",
            "book",
            "clock",
            "vase",
            "scissors",
            "teddy bear",
            "hair drier",
            "toothbrush",
        ]
    )
)


class Detection(NamedTuple):
    frame: float
    timedelta: datetime.timedelta
    xmin: float
    ymin: float
    xmax: float
    ymax: float
    xhalf: float
    yhalf: float
    confidence: float
    class_id: int
    class_name: str

    @staticmethod
    def from_xyxy(
        xyxy: Any, current_frame: float, timedelta: datetime.timedelta
    ) -> "Detection":
        return _detection_from_xyxy(xyxy, current_frame, timedelta)


def _detection_from_xyxy(
    xyxy: Any, current_frame: float, timedelta: datetime.timedelta
) -> Detection:
    xmin, ymin, xmax, ymax, confidence, class_id = xyxy
    xmin = float(xmin)
    ymin = float(ymin)
    xmax = float(xmax)
    ymax = float(ymax)
    xhalf = xmax - ((xmax - xmin) / 2)
    yhalf = ymax - ((ymax - ymin) / 2)
    confidence = float(confidence)
    class_id = int(class_id)
    class_name = class_name_by_class_id.get(class_id)

    return Detection(
        frame=current_frame,
        timedelta=timedelta,
        xmin=xmin,
        ymin=ymin,
        xmax=xmax,
        ymax=ymax,
        xhalf=xhalf,
        yhalf=yhalf,
        confidence=confidence,
        class_id=class_id,
        class_name=class_name,
    )


class DetectedFrame(NamedTuple):
    frame: float
    timedelta: datetime.timedelta
    detections: List[Detection]
