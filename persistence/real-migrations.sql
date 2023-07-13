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

CREATE SCHEMA public;

ALTER SCHEMA public OWNER TO postgres;

COMMENT ON SCHEMA public IS 'standard public schema';

SET default_tablespace = '';

SET default_table_access_method = heap;

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


COPY public.camera (id, uuid, name, stream_url, external_id) FROM stdin;
2	cd056389-b0b0-4978-9167-68c93e59f53d	FrontDoor	rtsp://192.168.137.32:554/Streaming/Channels/101/	102
1	3830e9a5-673d-4e7f-ae9b-afa9aeb439ab	Driveway	rtsp://192.168.137.31:554/Streaming/Channels/101/	101
3	ba9a4013-9c25-4bd0-931f-6eaf61e7369f	SideGate	rtsp://192.168.137.33:554/Streaming/Channels/101/	103
\.

COPY public.event (id, uuid, start_timestamp, end_timestamp, is_segment, high_quality_video_id, high_quality_image_id,
                   low_quality_video_id, low_quality_image_id, source_camera_id, needs_object_processing) FROM stdin;
1	782a3f31-f044-466b-9240-df182eea10e2	2023-05-14 13:15:03+00	2023-05-14 13:20:03+00	t	1	2	2	1	2	t
2	9053fa60-8f05-4b32-9b31-365d88081561	2023-05-14 13:20:32+00	2023-05-14 13:20:43+00	f	3	4	4	3	1	t
3	1c330cd7-16af-4234-ba93-a7bc29ad7d81	2023-05-14 13:15:00+00	2023-05-14 13:20:00+00	t	5	6	6	5	1	t
4	f52ab3b7-7aac-4bd8-8b97-852a0ba18338	2023-05-14 13:15:04+00	2023-05-14 13:20:04+00	t	7	8	8	7	3	t
5	00a120e3-990d-4880-b598-27b4e11dace9	2023-05-14 13:20:03+00	2023-05-14 13:25:03+00	t	9	10	10	9	2	t
6	8910357a-00ff-46be-ad38-b2352191a7ec	2023-05-14 13:20:00+00	2023-05-14 13:25:00+00	t	11	12	12	11	1	t
7	0960e413-9928-483c-bd9d-4ae58adef79b	2023-05-14 13:20:04+00	2023-05-14 13:25:04+00	t	13	14	14	13	3	t
\.

COPY public.image (id, uuid, "timestamp", size, is_high_quality, file_path, source_camera_id) FROM stdin;
1	bef0ac26-6548-40f0-a462-90fa80b2bb04	2023-05-14 13:15:03+00	0.032375	t	/srv/target_dir/segments/Segment_2023-05-14T21:15:03_FrontDoor__lowres.jpg	2
2	54528a19-5839-4cfd-a250-9117ae62dd24	2023-05-14 13:15:03+00	0.121444	t	/srv/target_dir/segments/Segment_2023-05-14T21:15:03_FrontDoor.jpg	2
3	0ce2ffc3-e6de-4e5d-ba60-8fc2703e2487	2023-05-14 13:20:32+00	0.153921	t	/srv/target_dir/events/Event_2023-05-14T21:20:31__101__Driveway__02__lowres.jpg	1
4	e5d67bae-f79e-4b58-b9be-34d40b039a44	2023-05-14 13:20:32+00	1.091943	t	/srv/target_dir/events/Event_2023-05-14T21:20:31__101__Driveway__02.jpg	1
5	9a6eae7c-d1e6-4b1a-ad40-5ed3864bccd2	2023-05-14 13:15:00+00	0.03441	t	/srv/target_dir/segments/Segment_2023-05-14T21:15:00_Driveway__lowres.jpg	1
6	3a0e10d7-7a34-46d1-b77b-a15b81d46022	2023-05-14 13:15:00+00	0.118712	t	/srv/target_dir/segments/Segment_2023-05-14T21:15:00_Driveway.jpg	1
7	edd63d0b-c035-4bb5-a596-85a909f936f4	2023-05-14 13:15:04+00	0.009595	t	/srv/target_dir/segments/Segment_2023-05-14T21:15:04_SideGate__lowres.jpg	3
8	1ffdb769-ee54-4fea-870a-e03d289b78e3	2023-05-14 13:15:04+00	0.030757	t	/srv/target_dir/segments/Segment_2023-05-14T21:15:04_SideGate.jpg	3
9	16743295-8cdb-4d2b-9455-c34bef681e6e	2023-05-14 13:20:03+00	0.03251	t	/srv/target_dir/segments/Segment_2023-05-14T21:20:03_FrontDoor__lowres.jpg	2
10	b9a41087-2356-4480-8d09-9e43b190ae22	2023-05-14 13:20:03+00	0.123442	t	/srv/target_dir/segments/Segment_2023-05-14T21:20:03_FrontDoor.jpg	2
11	ded49b97-2dbe-4126-99fb-cdfcc5613cff	2023-05-14 13:20:00+00	0.034408	t	/srv/target_dir/segments/Segment_2023-05-14T21:20:00_Driveway__lowres.jpg	1
12	227c79fb-af6b-4b56-a10e-d70c2b7a3493	2023-05-14 13:20:00+00	0.118564	t	/srv/target_dir/segments/Segment_2023-05-14T21:20:00_Driveway.jpg	1
13	10c1dd74-ca2e-4030-8c9b-5601ab88b89f	2023-05-14 13:20:04+00	0.009652	t	/srv/target_dir/segments/Segment_2023-05-14T21:20:04_SideGate__lowres.jpg	3
14	accfd407-bc0c-4791-8ed4-03de052f1897	2023-05-14 13:20:04+00	0.03076	t	/srv/target_dir/segments/Segment_2023-05-14T21:20:04_SideGate.jpg	3
\.

COPY public.object (id, uuid, start_timestamp, end_timestamp, detected_class_id, detected_class_name, tracked_object_id,
                    event_id, processed_video_id) FROM stdin;
\.

COPY public.video (id, uuid, start_timestamp, end_timestamp, size, is_high_quality, file_path,
                   source_camera_id) FROM stdin;
1	5ee40fc3-abb5-4e72-a557-40e090d72481	2023-05-14 13:15:03+00	2023-05-14 13:20:03+00	76.850061	t	/srv/target_dir/segments/Segment_2023-05-14T21:15:03_FrontDoor.mp4	2
2	42d28d82-567e-4989-91f1-0aaa5bf45b11	2023-05-14 13:15:03+00	2023-05-14 13:20:03+00	9.414657	t	/srv/target_dir/segments/Segment_2023-05-14T21:15:03_FrontDoor__lowres.mp4	2
3	d84af303-5fe7-4db1-8e8a-b82c30eb8786	2023-05-14 13:20:32+00	2023-05-14 13:20:43+00	6.419817	t	/srv/target_dir/events/Event_2023-05-14T21:20:27__101__Driveway__02.mp4	1
4	fb351037-7368-4986-9059-7aa4c5e648f2	2023-05-14 13:20:32+00	2023-05-14 13:20:43+00	1.51983	t	/srv/target_dir/events/Event_2023-05-14T21:20:27__101__Driveway__02__lowres.mp4	1
5	cedc30e5-8551-4f6c-88c2-35141eb466fb	2023-05-14 13:15:00+00	2023-05-14 13:20:00+00	76.905214	t	/srv/target_dir/segments/Segment_2023-05-14T21:15:00_Driveway.mp4	1
6	da416509-5df5-4800-ba69-6d2fb16f6488	2023-05-14 13:15:00+00	2023-05-14 13:20:00+00	10.934064	t	/srv/target_dir/segments/Segment_2023-05-14T21:15:00_Driveway__lowres.mp4	1
7	c0af59d5-c736-4dff-b8cf-8fd8fcb4745f	2023-05-14 13:15:04+00	2023-05-14 13:20:04+00	38.061661	t	/srv/target_dir/segments/Segment_2023-05-14T21:15:04_SideGate.mp4	3
8	e238f5a7-4f35-455b-be85-c35a046d505c	2023-05-14 13:15:04+00	2023-05-14 13:20:04+00	2.847458	t	/srv/target_dir/segments/Segment_2023-05-14T21:15:04_SideGate__lowres.mp4	3
9	23deef1d-d85e-4dc0-ba25-353d99ce183f	2023-05-14 13:20:03+00	2023-05-14 13:25:03+00	76.927582	t	/srv/target_dir/segments/Segment_2023-05-14T21:20:03_FrontDoor.mp4	2
10	f1ca424f-69c1-457d-b1eb-ce337b9929e6	2023-05-14 13:20:03+00	2023-05-14 13:25:03+00	9.42343	t	/srv/target_dir/segments/Segment_2023-05-14T21:20:03_FrontDoor__lowres.mp4	2
11	f6f1358f-0d10-4b9c-9557-498c2fc75705	2023-05-14 13:20:00+00	2023-05-14 13:25:00+00	76.909974	t	/srv/target_dir/segments/Segment_2023-05-14T21:20:00_Driveway.mp4	1
12	4c53127d-e8f3-477c-95cd-fe3c0a5b4bd5	2023-05-14 13:20:00+00	2023-05-14 13:25:00+00	11.394186	t	/srv/target_dir/segments/Segment_2023-05-14T21:20:00_Driveway__lowres.mp4	1
13	a99326be-d398-4fdc-84dd-6c1af76e24d2	2023-05-14 13:20:04+00	2023-05-14 13:25:04+00	38.203953	t	/srv/target_dir/segments/Segment_2023-05-14T21:20:04_SideGate.mp4	3
14	885f65eb-28e7-41b7-9138-4557c5079a9c	2023-05-14 13:20:04+00	2023-05-14 13:25:04+00	2.857351	t	/srv/target_dir/segments/Segment_2023-05-14T21:20:04_SideGate__lowres.mp4	3
\.

SELECT pg_catalog.setval('public.camera_id_seq', 3, true);


SELECT pg_catalog.setval('public.event_id_seq', 7, true);


SELECT pg_catalog.setval('public.image_id_seq', 14, true);


SELECT pg_catalog.setval('public.object_id_seq', 1, false);


SELECT pg_catalog.setval('public.video_id_seq', 14, true);


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
