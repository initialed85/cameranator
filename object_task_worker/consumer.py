from typing import Optional

from amqpy import Connection, Message, Timeout, Channel
from orjson import loads

from object_task_worker.object_tracker import (
    ObjectTracker,
    DEFAULT_CONF_THRESHOLD,
    DEFAULT_IOU_THRESHOLD,
    DEFAULT_IMAGE_SIZE,
    DEFAULT_STRIDE_FRAMES,
)

_AMQP_IDENTIFIER = "object_tasks"
_CONNECT_TIMEOUT = 5
_HEARTBEAT = 5


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
        self._ch: Optional[Channel] = None

        print("created object tracker..")
        self._object_tracker = ObjectTracker(
            model_name="yolov7.pt",
            device="cuda",
            tracking_mode="centroid",
            conf_threshold=DEFAULT_CONF_THRESHOLD,
            iou_threshold=DEFAULT_IOU_THRESHOLD,
            img_size=DEFAULT_IMAGE_SIZE,
            stride_frames=DEFAULT_STRIDE_FRAMES,
        )
        print("created {}".format(repr(self._object_tracker)))

    def _handler(self, message: Message):
        # event = {
        #     "uuid": "9c9c7858-230f-41f4-9fe0-1450c53c666a",
        #     "start_timestamp": "2023-01-31T09:01:48Z",
        #     "end_timestamp": "2023-01-31T09:02:03Z",
        #     "high_quality_video": {
        #         "file_path": "/srv/target_dir/events/Event_2023-01-31T17:01:44__102__FrontDoor__2710.mp4"
        #     },
        # }
        event = loads(message.body)
        print(f"received event={repr(event)}")

        input_path = event.get("high_quality_video", {}).get("file_path")
        if not input_path:
            raise ValueError(
                "failed to find high_quality_video::file_path in event={}".format(event)
            )

        processed_video = self._object_tracker(input_path)

        _ = processed_video.to_json()

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

        print("opening channel and declaring exchange and queue")
        self._ch = self._conn.channel()
        self._ch.exchange_declare(_AMQP_IDENTIFIER, "direct")
        self._ch.queue_declare(_AMQP_IDENTIFIER)
        self._ch.queue_bind(
            _AMQP_IDENTIFIER,
            exchange=_AMQP_IDENTIFIER,
            routing_key=_AMQP_IDENTIFIER,
        )

        self._ch.basic_consume(_AMQP_IDENTIFIER, callback=self._handler)

        print("waiting for events...")

        while 1:
            try:
                self._conn.drain_events(timeout=1)
            except Timeout:
                pass
            except KeyboardInterrupt:
                break

    def close(self):
        if self._ch:
            self._ch.close()

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
