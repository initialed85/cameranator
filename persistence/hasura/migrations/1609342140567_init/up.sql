SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

ALTER SCHEMA public OWNER TO postgres;

CREATE TABLE public.camera
(
    id          integer                        NOT NULL,
    uuid        uuid DEFAULT gen_random_uuid() NOT NULL,
    name        text                           NOT NULL,
    stream_url  text                           NOT NULL,
    external_id text
);

ALTER TABLE public.camera
    OWNER TO postgres;

CREATE SEQUENCE public.camera_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE public.camera_id_seq
    OWNER TO postgres;

ALTER SEQUENCE public.camera_id_seq OWNED BY public.camera.id;

CREATE TABLE public.event
(
    id                      integer                           NOT NULL,
    uuid                    uuid    DEFAULT gen_random_uuid() NOT NULL,
    start_timestamp         timestamp with time zone          NOT NULL,
    end_timestamp           timestamp with time zone          NOT NULL,
    is_segment              boolean DEFAULT false             NOT NULL,
    high_quality_video_id   integer                           NOT NULL,
    high_quality_image_id   integer                           NOT NULL,
    low_quality_video_id    integer,
    low_quality_image_id    integer,
    source_camera_id        integer                           NOT NULL,
    needs_object_processing boolean DEFAULT true              NOT NULL
);

ALTER TABLE public.event
    OWNER TO postgres;

CREATE SEQUENCE public.event_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE public.event_id_seq
    OWNER TO postgres;

ALTER SEQUENCE public.event_id_seq OWNED BY public.event.id;

CREATE TABLE public.image
(
    id               integer                        NOT NULL,
    uuid             uuid DEFAULT gen_random_uuid() NOT NULL,
    "timestamp"      timestamp with time zone       NOT NULL,
    size             double precision               NOT NULL,
    is_high_quality  boolean                        NOT NULL,
    file_path        text                           NOT NULL,
    source_camera_id integer                        NOT NULL
);

ALTER TABLE public.image
    OWNER TO postgres;

CREATE SEQUENCE public.image_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE public.image_id_seq
    OWNER TO postgres;

ALTER SEQUENCE public.image_id_seq OWNED BY public.image.id;

CREATE TABLE public.object
(
    id                  integer                        NOT NULL,
    uuid                uuid DEFAULT gen_random_uuid() NOT NULL,
    start_timestamp     timestamp with time zone       NOT NULL,
    end_timestamp       timestamp with time zone       NOT NULL,
    detected_class_id   integer                        NOT NULL,
    detected_class_name text                           NOT NULL,
    tracked_object_id   integer                        NOT NULL,
    event_id            integer                        NOT NULL,
    processed_video_id  integer                        NOT NULL
);

ALTER TABLE public.object
    OWNER TO postgres;

CREATE SEQUENCE public.object_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE public.object_id_seq
    OWNER TO postgres;

ALTER SEQUENCE public.object_id_seq OWNED BY public.object.id;

CREATE TABLE public.video
(
    id               integer                        NOT NULL,
    uuid             uuid DEFAULT gen_random_uuid() NOT NULL,
    start_timestamp  timestamp with time zone       NOT NULL,
    end_timestamp    timestamp with time zone       NOT NULL,
    size             double precision               NOT NULL,
    is_high_quality  boolean                        NOT NULL,
    file_path        text                           NOT NULL,
    source_camera_id integer                        NOT NULL
);

ALTER TABLE public.video
    OWNER TO postgres;

CREATE SEQUENCE public.video_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE public.video_id_seq
    OWNER TO postgres;

ALTER SEQUENCE public.video_id_seq OWNED BY public.video.id;

ALTER TABLE ONLY public.camera
    ALTER COLUMN id SET DEFAULT nextval('public.camera_id_seq'::regclass);

ALTER TABLE ONLY public.event
    ALTER COLUMN id SET DEFAULT nextval('public.event_id_seq'::regclass);

ALTER TABLE ONLY public.image
    ALTER COLUMN id SET DEFAULT nextval('public.image_id_seq'::regclass);

ALTER TABLE ONLY public.object
    ALTER COLUMN id SET DEFAULT nextval('public.object_id_seq'::regclass);

ALTER TABLE ONLY public.video
    ALTER COLUMN id SET DEFAULT nextval('public.video_id_seq'::regclass);

ALTER TABLE ONLY public.camera
    ADD CONSTRAINT camera_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.camera
    ADD CONSTRAINT camera_uuid_key UNIQUE (uuid);

ALTER TABLE ONLY public.event
    ADD CONSTRAINT event_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.event
    ADD CONSTRAINT event_uuid_key UNIQUE (uuid);

ALTER TABLE ONLY public.image
    ADD CONSTRAINT image_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.image
    ADD CONSTRAINT image_uuid_key UNIQUE (uuid);

ALTER TABLE ONLY public.object
    ADD CONSTRAINT object_uuid_key UNIQUE (uuid);

ALTER TABLE ONLY public.video
    ADD CONSTRAINT video_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.video
    ADD CONSTRAINT video_uuid_key UNIQUE (uuid);

CREATE INDEX event_is_segment_start_timestamp_end_timestamp_source_camer_idx ON public.event USING btree (is_segment, start_timestamp, end_timestamp, source_camera_id);

CREATE INDEX event_start_timestamp_idx ON public.event USING btree (start_timestamp);

CREATE INDEX object_detected_class_id_idx ON public.object USING btree (detected_class_id);

CREATE INDEX object_event_id_idx ON public.object USING btree (event_id);

CREATE INDEX object_start_timestamp_end_timestamp_idx ON public.object USING btree (start_timestamp, end_timestamp);

CREATE INDEX object_start_timestamp_idx ON public.object USING btree (start_timestamp);

ALTER TABLE ONLY public.event
    ADD CONSTRAINT event_high_quality_image_id_fkey FOREIGN KEY (high_quality_image_id) REFERENCES public.image (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE ONLY public.event
    ADD CONSTRAINT event_high_quality_video_id_fkey FOREIGN KEY (high_quality_video_id) REFERENCES public.video (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE ONLY public.event
    ADD CONSTRAINT event_low_quality_image_id_fkey FOREIGN KEY (low_quality_image_id) REFERENCES public.image (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE ONLY public.event
    ADD CONSTRAINT event_low_quality_video_id_fkey FOREIGN KEY (low_quality_video_id) REFERENCES public.video (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE ONLY public.event
    ADD CONSTRAINT event_source_camera_fkey FOREIGN KEY (source_camera_id) REFERENCES public.camera (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE ONLY public.image
    ADD CONSTRAINT image_source_camera_id_fkey FOREIGN KEY (source_camera_id) REFERENCES public.camera (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE ONLY public.object
    ADD CONSTRAINT object_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.event (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE ONLY public.object
    ADD CONSTRAINT object_processed_video_id_fkey FOREIGN KEY (processed_video_id) REFERENCES public.video (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE ONLY public.video
    ADD CONSTRAINT video_source_camera_id_fkey FOREIGN KEY (source_camera_id) REFERENCES public.camera (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

REVOKE USAGE ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO PUBLIC;
