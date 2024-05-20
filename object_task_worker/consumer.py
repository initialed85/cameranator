import os
import copy
import time
import traceback
from types import SimpleNamespace
import pika
from pika import BlockingConnection
from pika.channel import Channel
from pika.exceptions import ChannelClosed
from pika.spec import Basic, BasicProperties
import psycopg2
import datetime
from threading import Thread, Event
from queue import Queue, Empty

from typing import Optional, Dict, cast, List, Tuple
from dateutil.parser import parse
from orjson import loads

from object_task_worker.helpers import (
    get_query_for_repeaters,
    get_repeaters_for_detection,
)
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
_CONNECT_TIMEOUT = 5
_HEARTBEAT = 0

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

        self._conn: Optional[BlockingConnection] = None
        self._consume_ch: Optional[Channel] = None

        self._queue = Queue(maxsize=65536)
        self._stop_event = Event()

        self._thread_1 = Thread(
            target=self._handle_write_detections_to_db,
            daemon=False,
        )
        self._thread_1.start()

    def __del__(self):
        self._stop_event.set()

        try:
            print("trying to join background threads...")
            self._thread_1.join()
            print("joined.")
        except Exception:
            pass

    def _handle_write_detections_to_db(self):
        last_iteration = time.time() - 10.0

        while not self._stop_event.is_set():
            sleep = time.time() - last_iteration
            if sleep > 0:
                time.sleep(sleep)

            last_iteration = time.time()

            with psycopg2.connect(_DSN) as conn:
                with conn.cursor() as cur:
                    conn.set_session(autocommit=True)

                    tasks: List[
                        Tuple[
                            DetectionContext,
                            Dict[int, str],
                            datetime.datetime,
                            int,
                            int,
                        ]
                    ] = []

                    detections = 0
                    event_ids = []
                    repeaters = []

                    while (
                        not self._stop_event.is_set()
                        and time.time() - last_iteration < 10.0
                    ):
                        try:
                            (
                                detection_context,
                                name_by_class_id,
                                start_timestamp,
                                camera_id,
                                event_id,
                            ) = self._queue.get(timeout=1)
                        except Empty:
                            print("no tasks from queue...")
                            continue

                        detection_context = cast(DetectionContext, detection_context)
                        name_by_class_id = cast(Dict[int, str], name_by_class_id)
                        start_timestamp = cast(datetime.datetime, start_timestamp)
                        camera_id = cast(int, camera_id)
                        event_id = cast(int, event_id)

                        detections += len(detection_context.centroid_detections or [])
                        detections += len(detection_context.bbox_detections or [])

                        if event_id not in event_ids:
                            event_ids.append(event_id)

                        tasks.append(
                            (
                                detection_context,
                                name_by_class_id,
                                start_timestamp,
                                camera_id,
                                event_id,
                            )
                        )

                    if not len(tasks):
                        print(f"nothing to insert for this batch run")
                        continue

                    print(
                        f"buildin queries for insert of {detections} detections for event_ids={event_ids}"
                    )

                    for (
                        detection_context,
                        name_by_class_id,
                        start_timestamp,
                        camera_id,
                        event_id,
                    ) in tasks:
                        try:
                            repeaters.extend(
                                get_repeaters_for_detection(
                                    detection_context=detection_context,
                                    name_by_class_id=name_by_class_id,
                                    start_timestamp=start_timestamp,
                                    camera_id=camera_id,
                                    event_id=event_id,
                                )
                            )
                        except Exception:
                            traceback.print_exc()
                            continue

                    query = get_query_for_repeaters(repeaters)

                    print(
                        f"executing insert of {detections} detections for event_ids={event_ids}"
                    )

                    if not query:
                        continue

                    try:
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
        detections = max(
            [
                len(detection_context.centroid_detections),
                len(detection_context.bbox_detections),
            ]
        )

        print(f"queued {detections} for event_id={event_id} for writing to db")

        self._queue.put(
            (
                detection_context,
                name_by_class_id,
                start_timestamp,
                camera_id,
                event_id,
            )
        )

    def _actual_handler(self, message: SimpleNamespace):
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
            detection_context = copy.deepcopy(detection_context)
            name_by_class_id = copy.deepcopy(name_by_class_id)

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

        try:
            size = os.stat(processed_video.output_path).st_size / 1024 / 1024
        except Exception:
            size = 0
            with open(processed_video.output_path, "wb") as f:
                f.write(b"")

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

                for detected_object in processed_video.detected_objects or []:
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

    def _handler(
        self,
        channel: Channel,
        method: Basic.Deliver,
        properties: BasicProperties,
        body: bytes,
    ):
        message = SimpleNamespace(
            channel=channel,
            method=method,
            properties=properties,
            body=body,
        )

        try:
            print(f"handling message={repr(message)}...")
            self._actual_handler(message)
            print(f"successfully handled message={repr(message)}, acking...")
            channel.basic_ack(method.delivery_tag)
            print(f"acked message={repr(message)}.")
        except Exception as e:
            traceback.print_exc()

            print(
                f"failed to handle message={repr(message)}; e={repr(e)}, rejecting..."
            )

            try:
                channel.basic_reject(method.delivery_tag)
                print(f"rejected message={repr(message)}.")
            except Exception as e:
                traceback.print_exc()

                print(
                    f"failed to handle rejection of message={repr(message)}; e={repr(e)}, crashing out..."
                )

                self._stop_event.set()
                raise SystemExit(e) from e

    def open(self):
        print("connecting to {}:{}".format(self._host, self._port))

        credentials = pika.PlainCredentials(self._userid, self._password)
        parameters = pika.ConnectionParameters(self._host, credentials=credentials)
        self._conn = pika.BlockingConnection(parameters)

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

        self._consume_ch.basic_qos(
            prefetch_count=1,
        )

        self._consume_ch.basic_consume(
            _AMQP_IDENTIFIER,
            on_message_callback=self._handler,
        )

        print("consuming...")

        try:
            self._consume_ch.start_consuming()
        except ChannelClosed:
            traceback.print_exc()

            print(f"failed to consume; e={repr(e)}, crashing out...")

            self._stop_event.set()
            raise SystemExit(e) from e
        except KeyboardInterrupt:
            self._consume_ch.stop_consuming()

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
