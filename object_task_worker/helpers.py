import numpy as np
import datetime
from norfair import Detection
from typing import List, Optional, NamedTuple, Tuple, Dict


class DetectionContext(NamedTuple):
    centroid_detections: List[Detection]
    bbox_detections: List[Detection]
    timedelta: Optional[datetime.timedelta]


_VALUES_REPEATER = """(
    '{}',
    {},
    '{}',
    {},
    {},
    {},
    {},
    {}
)
"""

_INSERT_DETECTION_QUERY = """
INSERT INTO detection (
    timestamp,
    class_id,
    class_name,
    score,
    centroid,
    bounding_box,
    camera_id,
    event_id
) VALUES
{};
"""


def get_query_for_detection(
    detection_context: DetectionContext,
    name_by_class_id: Dict[int, str],
    start_timestamp: datetime.datetime,
    camera_id: int,
    event_id: int,
) -> Optional[str]:
    # TODO: use detection_context.*_detections.[].scores.mean() for confidence

    repeaters: List[str] = []

    for i, centroid_detection in enumerate(detection_context.centroid_detections or []):
        bbox_detection = detection_context.bbox_detections[i]

        score = (centroid_detection.scores.mean() + bbox_detection.scores.mean()) / 2

        tlx, tly = bbox_detection.points[0]
        brx, bry = bbox_detection.points[1]

        parts = [
            f"{tlx} {tly}",  # top left
            f"{brx} {tly}",  # top right
            f"{brx} {bry}",  # bottom right
            f"{tlx} {bry}",  # bottom left
            f"{tlx} {tly}",  # top left again (to close polygon)
        ]

        line_string = ", ".join(parts)

        repeaters.append(
            _VALUES_REPEATER.format(
                start_timestamp + detection_context.timedelta,
                centroid_detection.label,
                name_by_class_id.get(centroid_detection.label) or "__unknown__",
                score,
                f"ST_MakePoint({centroid_detection.points[0][0]}, {centroid_detection.points[0][1]})::Point",
                f"ST_MakePolygon(ST_GeomFromText('LINESTRING({line_string})'))::Polygon",
                camera_id,
                event_id,
            )
        )

    if not repeaters:
        return None

    return _INSERT_DETECTION_QUERY.format(", ".join(repeaters))
