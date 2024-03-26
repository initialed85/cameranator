import numpy as np
import datetime
from norfair import Detection
from typing import List, Optional, NamedTuple, Tuple, Dict


class DetectionContext(NamedTuple):
    centroid_detections: List[Detection]
    bbox_detections: List[Detection]
    dominant_colours: List[np.array]
    timedelta: Optional[datetime.timedelta]


_VALUES_REPEATER = """(
    '{}',
    {},
    '{}',
    {},
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
    colour,
    camera_id,
    event_id
) VALUES
{};
"""


def get_repeaters_for_detection(
    detection_context: DetectionContext,
    name_by_class_id: Dict[int, str],
    start_timestamp: datetime.datetime,
    camera_id: int,
    event_id: int,
) -> List[str]:
    repeaters: List[str] = []

    for i, centroid_detection in enumerate(detection_context.centroid_detections or []):
        bbox_detection = detection_context.bbox_detections[i]
        dominant_colour = detection_context.dominant_colours[i]

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
                f"ST_MakePoint({dominant_colour[0]}, {dominant_colour[1]}, {dominant_colour[2]})::geometry(pointz)",
                camera_id,
                event_id,
            )
        )

    return repeaters


def get_query_for_repeaters(repeaters: List[str]) -> Optional[str]:
    if not repeaters:
        return None

    return _INSERT_DETECTION_QUERY.format(", ".join(repeaters))
