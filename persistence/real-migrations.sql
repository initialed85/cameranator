TRUNCATE TABLE camera CASCADE;

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

CREATE EXTENSION IF NOT EXISTS postgis_raster;

SET
    postgis.gdal_enabled_drivers = 'ENABLE_ALL';

--
-- camera
--
CREATE TABLE
    public.camera (id bigint NOT NULL, name text NOT NULL, stream_url text NOT NULL);

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

ALTER TABLE ONLY public.camera
ADD CONSTRAINT camera_pkey PRIMARY KEY (id);

--
-- event
--
CREATE TABLE
    public.event (
        id bigint NOT NULL,
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

ALTER TABLE ONLY public.event
ADD CONSTRAINT event_pkey PRIMARY KEY (id);

--
-- video
--
CREATE TABLE
    public.video (
        id bigint NOT NULL,
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

ALTER TABLE ONLY public.video
ADD CONSTRAINT video_pkey PRIMARY KEY (id);

--
-- image
--
CREATE TABLE
    public.image (
        id bigint NOT NULL,
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

ALTER TABLE ONLY public.image
ADD CONSTRAINT image_pkey PRIMARY KEY (id);

--
-- object
--
CREATE TABLE
    public.object (
        id bigint NOT NULL,
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

ALTER TABLE ONLY public.object
ADD CONSTRAINT object_pkey PRIMARY KEY (id);

--
-- detection
--
CREATE TABLE
    public.detection (
        id bigint NOT NULL,
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

ALTER TABLE ONLY public.detection
ADD CONSTRAINT detection_pkey PRIMARY KEY (id);

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
-- seed data
--
INSERT INTO
    public.camera (id, name, stream_url)
VALUES
    (1, 'Driveway', 'rtsp://192.168.137.31:554/Streaming/Channels/101/'),
    (2, 'FrontDoor', 'rtsp://192.168.137.32:554/Streaming/Channels/101/'),
    (3, 'SideGate', 'rtsp://192.168.137.33:554/Streaming/Channels/101/') ON CONFLICT (id)
DO NOTHING;

SELECT
    pg_catalog.setval ('public.camera_id_seq', 3, true);
