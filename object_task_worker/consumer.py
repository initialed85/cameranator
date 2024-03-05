import os
from typing import Optional

import psycopg2
from amqpy import Connection, Message, Timeout, Channel
from dateutil.parser import parse
from orjson import loads

from object_task_worker.object_tracker import (
    ObjectTracker,
    DEFAULT_CONF_THRESHOLD,
    DEFAULT_IOU_THRESHOLD,
    DEFAULT_IMAGE_SIZE,
    DEFAULT_STRIDE_FRAMES,
    ProcessedVideo,
)

_AMQP_IDENTIFIER = "object_tasks"
_CONNECT_TIMEOUT = 60
_HEARTBEAT = 60

_INSERT_VIDEO_QUERY = """
INSERT INTO video (
    start_timestamp, end_timestamp, size, is_high_quality, file_path, source_camera_id
) VALUES (
    %s,
    %s,
    %s,
    true,
    %s,
    %s
)
RETURNING id;
"""

_INSERT_OBJECT_QUERY = """
INSERT INTO object (
    start_timestamp, end_timestamp, detected_class_id, detected_class_name, tracked_object_id, event_id, processed_video_id
) VALUES (
    %s,
    %s,
    %s,
    %s,
    %s,
    %s,
    %s
);
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
        # self._produce_ch: Optional[Channel] = None

    def _actual_handler(self, message: Message):
        # example event
        _ = {
            "id": 69,
            "high_quality_video": {
                "file_path": "/srv/target_dir/events/Event_2023-01-31T17:01:44__102__FrontDoor__2710.mp4",
                "source_camera_id": 3,
                "start_timestamp": "2023-01-31T09:01:48Z",
                "end_timestamp": "2023-01-31T09:02:03Z",
            },
        }

        event = loads(message.body)
        print(f"received event={repr(event)}")

        event_id = event.get("id")
        file_path = event.get("high_quality_video", {}).get("file_path")
        source_camera_id = event.get("high_quality_video", {}).get("source_camera_id")
        start_timestamp = event.get("high_quality_video", {}).get("start_timestamp")
        end_timestamp = event.get("high_quality_video", {}).get("end_timestamp")

        if not all(
            [event_id, file_path, source_camera_id, start_timestamp, end_timestamp]
        ):
            raise ValueError("unexpectedly Falsey field in event={}".format(event))

        start_timestamp = parse(start_timestamp)
        end_timestamp = parse(end_timestamp)

        print("creating object tracker..")
        object_tracker = ObjectTracker(
            model_name="yolov7.pt",
            device="cuda",
            tracking_mode="bbox",
            conf_threshold=DEFAULT_CONF_THRESHOLD,
            iou_threshold=DEFAULT_IOU_THRESHOLD,
            img_size=DEFAULT_IMAGE_SIZE,
            stride_frames=DEFAULT_STRIDE_FRAMES,
        )
        print("created {}".format(repr(object_tracker)))

        print("processing video...")
        processed_video: ProcessedVideo = object_tracker(
            file_path
        )  # slow, blocking call to do the processing

        size = os.stat(processed_video.output_path).st_size / 1024 / 1024

        with psycopg2.connect(_DSN) as conn:
            with conn.cursor() as cur:
                print("insert video row...")
                cur.execute(
                    _INSERT_VIDEO_QUERY,
                    (
                        start_timestamp.isoformat(),
                        end_timestamp.isoformat(),
                        size,
                        processed_video.output_path,
                        source_camera_id,
                    ),
                )

                processed_video_id = cur.fetchone()[0]

                for detected_object in processed_video.detected_objects:
                    print(
                        f"insert object row for object_id={detected_object.object_id}..."
                    )
                    cur.execute(
                        _INSERT_OBJECT_QUERY,
                        (
                            start_timestamp + detected_object.start_timedelta,
                            start_timestamp + detected_object.end_timedelta,
                            detected_object.class_id,
                            detected_object.class_name,
                            detected_object.object_id,
                            event_id,
                            processed_video_id,
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

        # print("opening produce channel and declaring exchange and queue")
        # self._produce_ch = self._conn.channel()

        # self._produce_ch.exchange_declare(
        #     _AMQP_IDENTIFIER,
        #     "direct",
        # )

        # self._produce_ch.queue_declare(
        #     _AMQP_IDENTIFIER,
        # )

        # self._produce_ch.queue_bind(
        #     _AMQP_IDENTIFIER,
        #     exchange=_AMQP_IDENTIFIER,
        #     routing_key=_AMQP_IDENTIFIER,
        # )

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

        # if self._produce_ch:
        #     self._produce_ch.close()

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
