import os
import traceback
import psycopg2
import datetime
from threading import Thread, Event
from queue import Queue, Empty

from typing import Optional, Dict, cast
from amqpy import Connection, Message, Timeout, Channel
from dateutil.parser import parse
from orjson import loads

from object_task_worker.helpers import get_query_for_detection
from object_task_worker.object_tracker import (
    DetectionContext,
    ObjectTracker,
    CONF_THRESHOLD,
    IOU_THRESHOLD,
    IMAGE_SIZE,
    STRIDE_FRAMES,
    ProcessedVideo,
)

_AMQP_IDENTIFIER = "object_tasks"
_CONNECT_TIMEOUT = 60
_HEARTBEAT = 60

_INSERT_VIDEO_QUERY = """
INSERT INTO video (
    start_timestamp,
    end_timestamp,
    size,
    file_path,
    camera_id
) VALUES (
    %s,
    %s,
    %s,
    %s,
    %s
)
RETURNING id;
"""

_INSERT_OBJECT_QUERY = """
INSERT INTO object (
    start_timestamp,
    end_timestamp,
    class_id,
    class_name,
    camera_id,
    event_id
) VALUES (
    %s,
    %s,
    %s,
    %s,
    %s,
    %s
);
"""

_UPDATE_EVENT_QUERY = """
UPDATE event SET
    processed_video_id = %s,
    status = 'needs tracking'
WHERE
    id = %s;
"""

_DB_HOST = (os.getenv("DB_HOST") or "").strip()
_DB_PORT = (os.getenv("DB_PORT") or "").strip()
_DB_USER = (os.getenv("DB_USER") or "").strip()
_DB_PASSWORD = (os.getenv("DB_PASSWORD") or "").strip()
_DB_NAME = (os.getenv("DB_NAME") or "").strip()

if not all([_DB_HOST, _DB_PASSWORD, _DB_PORT, _DB_USER]):
    raise ValueError(
        "one or more of DB_HOST, DB_PASSWORD, DB_PORT, DB_USER empty or unset"
    )

_DSN = f"dbname={_DB_NAME} user={_DB_USER} host={_DB_HOST} port={_DB_PORT} password={_DB_PASSWORD}"


class Consumer(object):
    def __init__(
        self,
        host: str,
        port: int,
        userid: str,
        password: str,
    ):
        self._host = host
        self._port = port
        self._userid = userid
        self._password = password

        self._conn: Optional[Connection] = None
        self._consume_ch: Optional[Channel] = None

        self._queue = Queue(maxsize=65536)
        self._stop_event = Event()
        self._thread = Thread(target=self._handle_write_detections_to_db, daemon=False)
        self._thread.start()

    def __del__(self):
        self._stop_event.set()

        try:
            print("trying to join background thread...")
            self._thread.join()
            print("joined.")
        except Exception:
            pass

    def _handle_write_detections_to_db(self):
        while not self._stop_event.is_set():
            try:
                (
                    detection_context,
                    name_by_class_id,
                    start_timestamp,
                    camera_id,
                    event_id,
                ) = self._queue.get(timeout=1)
            except Empty:
                continue

            detection_context = cast(DetectionContext, detection_context)
            name_by_class_id = cast(Dict[int, str], name_by_class_id)
            start_timestamp = cast(datetime.datetime, start_timestamp)
            camera_id = cast(int, camera_id)
            event_id = cast(int, event_id)

            try:
                query = get_query_for_detection(
                    detection_context=detection_context,
                    name_by_class_id=name_by_class_id,
                    start_timestamp=start_timestamp,
                    camera_id=camera_id,
                    event_id=event_id,
                )
            except Exception:
                traceback.print_exc()
                continue

            if not query:
                continue

            print(query)

            with psycopg2.connect(_DSN) as conn:
                with conn.cursor() as cur:
                    try:
                        detections = len(
                            detection_context.centroid_detections or []
                        ) + len(detection_context.bbox_detections or [])

                        print(
                            f"insert of {detections} detections for event_id={event_id}"
                        )

                        cur.execute(query)
                    except Exception:
                        traceback.print_exc()
                        continue

    def _write_detections_to_db(
        self,
        detection_context: DetectionContext,
        name_by_class_id: Dict[int, str],
        start_timestamp: datetime.datetime,
        camera_id: int,
        event_id: int,
    ):
        self._queue.put(
            (
                detection_context,
                name_by_class_id,
                start_timestamp,
                camera_id,
                event_id,
            )
        )

    def _actual_handler(self, message: Message):
        event = loads(message.body)
        print(f"received event={repr(event)}")

        event_id = event.get("id")
        file_path = event.get("original_video", {}).get("file_path")
        camera_id = event.get("original_video", {}).get("camera_id")
        start_timestamp = event.get("original_video", {}).get("start_timestamp")
        end_timestamp = event.get("original_video", {}).get("end_timestamp")

        if not all([event_id, file_path, camera_id, start_timestamp, end_timestamp]):
            raise ValueError("unexpectedly Falsey field in event={}".format(event))

        start_timestamp = parse(start_timestamp)
        end_timestamp = parse(end_timestamp)

        def write_detections_to_db(
            detection_context: DetectionContext,
            name_by_class_id: Dict[int, str],
        ):
            print(
                f"queueing insert of {max([len(detection_context.centroid_detections), len(detection_context.bbox_detections)])} detections for event_id={event_id}"
            )

            return self._write_detections_to_db(
                detection_context=detection_context,
                name_by_class_id=name_by_class_id,
                start_timestamp=start_timestamp,
                camera_id=camera_id,
                event_id=event_id,
            )

        print("creating object tracker..")
        object_tracker = ObjectTracker(
            model_name="yolov7.pt",
            device="cuda",
            tracking_mode="bbox",  # bbox or centroid
            conf_threshold=CONF_THRESHOLD,
            iou_threshold=IOU_THRESHOLD,
            img_size=IMAGE_SIZE,
            stride_frames=STRIDE_FRAMES,
            handle_detections=write_detections_to_db,
        )
        print(f"created {repr(object_tracker)}")

        print(f"processing video for file_path={file_path}...")
        processed_video: ProcessedVideo = object_tracker(
            file_path
        )  # slow, blocking call to do the processing

        size = os.stat(processed_video.output_path).st_size / 1024 / 1024

        with psycopg2.connect(_DSN) as conn:
            with conn.cursor() as cur:
                print(
                    f"insert video row for file_path={repr(processed_video.output_path)}..."
                )
                cur.execute(
                    _INSERT_VIDEO_QUERY,
                    (
                        start_timestamp.isoformat(),
                        end_timestamp.isoformat(),
                        size,
                        processed_video.output_path,
                        camera_id,
                    ),
                )
                processed_video_id = cur.fetchone()[0]

                print(f"update event row for event_id={event_id}...")
                cur.execute(
                    _UPDATE_EVENT_QUERY,
                    (
                        processed_video_id,
                        event_id,
                    ),
                )

                for detected_object in processed_video.detected_objects:
                    print(
                        f"insert object row for class_id={detected_object.class_id}, class_name={repr(detected_object.class_name)}..."
                    )
                    cur.execute(
                        _INSERT_OBJECT_QUERY,
                        (
                            start_timestamp + detected_object.start_timedelta,
                            start_timestamp + detected_object.end_timedelta,
                            detected_object.class_id,
                            detected_object.class_name,
                            camera_id,
                            event_id,
                        ),
                    )

        print("done.")

    def _handler(self, message: Message):
        try:
            print(f"handling message={repr(message)}...")
            self._actual_handler(message)
            print(f"successfully handled message={repr(message)}, acking...")
            message.ack()
            print(f"acked message={repr(message)}.")
        except Exception as e:
            print(
                f"failed to handle message={repr(message)}; e={repr(e)}, rejecting..."
            )
            message.reject(requeue=True)
            print(f"rejecrted message={repr(message)}.")

    def open(self):
        print("connecting to {}:{}".format(self._host, self._port))
        self._conn = Connection(
            host=self._host,
            port=self._port,
            userid=self._userid,
            password=self._password,
            connect_timeout=_CONNECT_TIMEOUT,
            heartbeat=_HEARTBEAT,
        )

        print("opening consume channel and declaring exchange and queue")
        self._consume_ch = self._conn.channel()

        self._consume_ch.exchange_declare(
            _AMQP_IDENTIFIER,
            "direct",
            durable=True,
            auto_delete=False,
        )

        self._consume_ch.queue_declare(
            _AMQP_IDENTIFIER,
            durable=True,
            auto_delete=False,
        )

        self._consume_ch.queue_bind(
            _AMQP_IDENTIFIER,
            exchange=_AMQP_IDENTIFIER,
            routing_key=_AMQP_IDENTIFIER,
        )

        self._consume_ch.basic_consume(
            _AMQP_IDENTIFIER,
            callback=self._handler,
        )

        print("waiting for events...")

        while 1:
            try:
                self._conn.drain_events(timeout=5)
            except Timeout:
                pass
            except KeyboardInterrupt:
                break

    def close(self):
        if self._consume_ch:
            self._consume_ch.close()

        if self._conn:
            self._conn.close()


def run(
    host,
    port,
    userid,
    password,
):
    consumer = Consumer(
        host=host,
        port=port,
        userid=userid,
        password=password,
    )

    try:
        consumer.open()
    finally:
        consumer.close()
