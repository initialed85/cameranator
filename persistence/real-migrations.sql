DROP TABLE IF EXISTS public.camera CASCADE;

DROP TABLE IF EXISTS public.video CASCADE;

DROP TABLE IF EXISTS public.image CASCADE;

DROP TABLE IF EXISTS public.object CASCADE;

DROP TABLE IF EXISTS public.detection CASCADE;

SET
    statement_timeout = 0;

SET
    lock_timeout = 0;

SET
    idle_in_transaction_session_timeout = 0;

SET
    client_encoding = 'UTF8';

SET
    standard_conforming_strings = on;

SELECT
    pg_catalog.set_config ('search_path', '', false);

SET
    check_function_bodies = false;

SET
    xmloption = content;

SET
    client_min_messages = warning;

SET
    row_security = off;

CREATE SCHEMA IF NOT EXISTS public;

ALTER SCHEMA public OWNER TO postgres;

COMMENT ON SCHEMA public IS 'standard public schema';

SET
    default_tablespace = '';

SET
    default_table_access_method = heap;

CREATE EXTENSION IF NOT EXISTS postgis SCHEMA public;

CREATE EXTENSION IF NOT EXISTS postgis_raster SCHEMA public;

SET
    postgis.gdal_enabled_drivers = 'ENABLE_ALL';

--
-- camera
--
CREATE TABLE
    public.camera (id bigint NOT NULL PRIMARY KEY, name text NOT NULL UNIQUE, stream_url text NOT NULL);

ALTER TABLE public.camera OWNER TO postgres;

CREATE SEQUENCE public.camera_id_seq AS bigint START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER TABLE public.camera_id_seq OWNER TO postgres;

ALTER SEQUENCE public.camera_id_seq OWNED BY public.camera.id;

ALTER TABLE ONLY public.camera
ALTER COLUMN id
SET DEFAULT nextval('public.camera_id_seq'::regclass);

SELECT
    pg_catalog.setval ('public.camera_id_seq', 1, true);

--
-- event
--
CREATE TABLE
    public.event (
        id bigint NOT NULL PRIMARY KEY,
        start_timestamp timestamp with time zone NOT NULL,
        end_timestamp timestamp with time zone NOT NULL,
        duration interval NOT NULL DEFAULT interval '0 seconds',
        original_video_id bigint NOT NULL,
        thumbnail_image_id bigint NOT NULL,
        processed_video_id bigint,
        source_camera_id bigint NOT NULL,
        status text DEFAULT true NOT NULL CHECK (status IN ('needs detection', 'detection underway', 'needs tracking', 'tracking underway', 'done'))
    );

ALTER TABLE public.event OWNER TO postgres;

CREATE SEQUENCE public.event_id_seq AS bigint START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER TABLE public.event_id_seq OWNER TO postgres;

ALTER SEQUENCE public.event_id_seq OWNED BY public.event.id;

ALTER TABLE ONLY public.event
ALTER COLUMN id
SET DEFAULT nextval('public.event_id_seq'::regclass);

SELECT
    pg_catalog.setval ('public.event_id_seq', 1, true);

--
-- video
--
CREATE TABLE
    public.video (
        id bigint NOT NULL PRIMARY KEY,
        start_timestamp timestamp with time zone NOT NULL,
        end_timestamp timestamp with time zone NOT NULL,
        duration interval NOT NULL DEFAULT interval '0 seconds',
        "size" double precision NOT NULL DEFAULT 0,
        file_path text NOT NULL,
        camera_id bigint NOT NULL,
        event_id bigint NULL
    );

ALTER TABLE public.video OWNER TO postgres;

CREATE SEQUENCE public.video_id_seq AS bigint START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER TABLE public.video_id_seq OWNER TO postgres;

ALTER SEQUENCE public.video_id_seq OWNED BY public.video.id;

ALTER TABLE ONLY public.video
ALTER COLUMN id
SET DEFAULT nextval('public.video_id_seq'::regclass);

SELECT
    pg_catalog.setval ('public.video_id_seq', 1, true);

--
-- image
--
CREATE TABLE
    public.image (
        id bigint NOT NULL PRIMARY KEY,
        "timestamp" timestamp with time zone NOT NULL,
        "size" double precision NOT NULL DEFAULT 0,
        file_path text NOT NULL,
        camera_id bigint NOT NULL,
        event_id bigint NULL
    );

ALTER TABLE public.image OWNER TO postgres;

CREATE SEQUENCE public.image_id_seq AS bigint START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER TABLE public.image_id_seq OWNER TO postgres;

ALTER SEQUENCE public.image_id_seq OWNED BY public.image.id;

ALTER TABLE ONLY public.image
ALTER COLUMN id
SET DEFAULT nextval('public.image_id_seq'::regclass);

SELECT
    pg_catalog.setval ('public.image_id_seq', 1, true);

--
-- object
--
CREATE TABLE
    public.object (
        id bigint NOT NULL PRIMARY KEY,
        start_timestamp timestamp with time zone NOT NULL,
        end_timestamp timestamp with time zone NOT NULL,
        class_id bigint NOT NULL,
        class_name text NOT NULL,
        camera_id bigint NOT NULL,
        event_id bigint NULL
    );

ALTER TABLE public.object OWNER TO postgres;

CREATE SEQUENCE public.object_id_seq AS bigint START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER TABLE public.object_id_seq OWNER TO postgres;

ALTER SEQUENCE public.object_id_seq OWNED BY public.object.id;

ALTER TABLE ONLY public.object
ALTER COLUMN id
SET DEFAULT nextval('public.object_id_seq'::regclass);

SELECT
    pg_catalog.setval ('public.object_id_seq', 1, false);

--
-- detection
--
CREATE TABLE
    public.detection (
        id bigint NOT NULL PRIMARY KEY,
        "timestamp" timestamp with time zone NOT NULL,
        class_id bigint NOT NULL,
        class_name text NOT NULL,
        score float NOT NULL,
        centroid Point NOT NULL,
        bounding_box Polygon NOT NULL,
        camera_id bigint NOT NULL,
        event_id bigint NULL,
        object_id bigint NULL
    );

ALTER TABLE public.detection OWNER TO postgres;

CREATE SEQUENCE public.detection_id_seq AS bigint START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER TABLE public.detection_id_seq OWNER TO postgres;

ALTER SEQUENCE public.detection_id_seq OWNED BY public.detection.id;

ALTER TABLE ONLY public.detection
ALTER COLUMN id
SET DEFAULT nextval('public.detection_id_seq'::regclass);

SELECT
    pg_catalog.setval ('public.detection_id_seq', 1, true);

--
-- foreign keys
--
ALTER TABLE ONLY public.event
ADD CONSTRAINT event_original_video_id_fkey FOREIGN KEY (original_video_id) REFERENCES public.video (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE ONLY public.event
ADD CONSTRAINT event_thumbnail_image_id_fkey FOREIGN KEY (thumbnail_image_id) REFERENCES public.image (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE ONLY public.event
ADD CONSTRAINT event_processed_video_id_fkey FOREIGN KEY (processed_video_id) REFERENCES public.video (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE ONLY public.event
ADD CONSTRAINT event_source_camera_id_fkey FOREIGN KEY (source_camera_id) REFERENCES public.camera (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE ONLY public.video
ADD CONSTRAINT video_camera_id_fkey FOREIGN KEY (camera_id) REFERENCES public.camera (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE ONLY public.video
ADD CONSTRAINT video_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.event (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE ONLY public.image
ADD CONSTRAINT image_camera_id_fkey FOREIGN KEY (camera_id) REFERENCES public.camera (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE ONLY public.image
ADD CONSTRAINT image_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.event (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE ONLY public.object
ADD CONSTRAINT objet_camera_id_fkey FOREIGN KEY (camera_id) REFERENCES public.camera (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE ONLY public.object
ADD CONSTRAINT objet_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.event (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE ONLY public.detection
ADD CONSTRAINT detection_camera_id_fkey FOREIGN KEY (camera_id) REFERENCES public.camera (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE ONLY public.detection
ADD CONSTRAINT detection_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.event (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE ONLY public.detection
ADD CONSTRAINT detection_object_id_fkey FOREIGN KEY (object_id) REFERENCES public.object (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

--
-- because my neanderthal brain cannot switch contexts from the naming schema we use at work
--
DROP VIEW IF EXISTS public.cameras;

CREATE VIEW
    public.cameras AS
SELECT
    *
FROM
    public.camera;

DROP VIEW IF EXISTS public.detections;

CREATE VIEW
    public.detections AS
SELECT
    *
FROM
    public.detection;

DROP VIEW IF EXISTS public.events;

CREATE VIEW
    public.events AS
SELECT
    *
FROM
    public.event;

DROP VIEW IF EXISTS public.images;

CREATE VIEW
    public.images AS
SELECT
    *
FROM
    public.image;

DROP VIEW IF EXISTS public.objects;

CREATE VIEW
    public.objects AS
SELECT
    *
FROM
    public.object;

DROP VIEW IF EXISTS public.videos;

CREATE VIEW
    public.videos AS
SELECT
    *
FROM
    public.video;

--
-- indexes
--
CREATE INDEX IF NOT EXISTS event_start_timestamp_idx ON public.event (start_timestamp);

CREATE INDEX IF NOT EXISTS event_end_timestamp_idx ON public.event (end_timestamp);

CREATE INDEX IF NOT EXISTS video_start_timestamp_idx ON public.video (start_timestamp);

CREATE INDEX IF NOT EXISTS video_end_timestamp_idx ON public.video (end_timestamp);

CREATE INDEX IF NOT EXISTS image_timestamp_idx ON public.image ("timestamp");

CREATE INDEX IF NOT EXISTS object_start_timestamp_idx ON public.object (start_timestamp);

CREATE INDEX IF NOT EXISTS object_end_timestamp_idx ON public.object (end_timestamp);

CREATE INDEX IF NOT EXISTS object_class_id_idx ON public.object (class_id);

CREATE INDEX IF NOT EXISTS object_class_name_idx ON public.object (class_name);

CREATE INDEX IF NOT EXISTS detection_timestamp_idx ON public.detection ("timestamp");

CREATE INDEX IF NOT EXISTS detection_class_id_idx ON public.detection (class_id);

CREATE INDEX IF NOT EXISTS detection_class_name_idx ON public.detection (class_name);

--
-- views
--
DROP VIEW IF EXISTS event_with_detection;

CREATE VIEW
    event_with_detection AS (
        WITH
            detections_1 AS (
                SELECT
                    d.event_id,
                    d.class_id,
                    d.class_name,
                    avg(d.score) AS score,
                    count(d.id) AS count,
                    avg(d.score) * count(d.id) AS weighted_score
                FROM
                    detection d
                GROUP BY
                    (event_id, class_id, class_name)
            ),
            events_1 AS (
                SELECT
                    e.*,
                    d.*
                FROM
                    events e
                    LEFT JOIN detections_1 d ON d.event_id = e.id
            )
        SELECT
            *
        FROM
            events_1 e
    );

--
-- seed data
--
INSERT INTO
    public.camera (id, name, stream_url)
VALUES
    (1, 'Driveway', 'rtsp://192.168.137.31:554/Streaming/Channels/101/'),
    (2, 'FrontDoor', 'rtsp://192.168.137.32:554/Streaming/Channels/101/'),
    (3, 'SideGate', 'rtsp://192.168.137.33:554/Streaming/Channels/101/') ON CONFLICT (id)
DO NOTHING;

WITH
    cte AS (
        SELECT
            max(id) AS max_id
        FROM
            public.camera
    )
SELECT
    pg_catalog.setval ('public.camera_id_seq', cte.max_id, true)
FROM
    cte;
