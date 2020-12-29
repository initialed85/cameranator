CREATE TABLE public.camera (
    id integer NOT NULL,
    uuid uuid DEFAULT public.gen_random_uuid() NOT NULL,
    name text NOT NULL,
    stream_url text NOT NULL
);
CREATE SEQUENCE public.camera_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.camera_id_seq OWNED BY public.camera.id;
CREATE TABLE public.event (
    id integer NOT NULL,
    uuid uuid DEFAULT public.gen_random_uuid() NOT NULL,
    start_timestamp timestamp with time zone NOT NULL,
    end_timestamp timestamp with time zone NOT NULL,
    is_processed boolean DEFAULT false NOT NULL,
    high_quality_video_id integer NOT NULL,
    high_quality_image_id integer NOT NULL,
    low_quality_video_id integer,
    low_quality_image_id integer,
    source_camera_id integer NOT NULL
);
CREATE SEQUENCE public.event_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.event_id_seq OWNED BY public.event.id;
CREATE TABLE public.image (
    id integer NOT NULL,
    uuid uuid DEFAULT public.gen_random_uuid() NOT NULL,
    "timestamp" timestamp with time zone NOT NULL,
    size double precision NOT NULL,
    is_high_quality boolean NOT NULL,
    file_path text NOT NULL,
    source_camera_id integer NOT NULL
);
CREATE SEQUENCE public.image_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.image_id_seq OWNED BY public.image.id;
CREATE TABLE public.video (
    id integer NOT NULL,
    uuid uuid DEFAULT public.gen_random_uuid() NOT NULL,
    start_timestamp timestamp without time zone NOT NULL,
    end_timestamp timestamp with time zone NOT NULL,
    size double precision NOT NULL,
    is_high_quality boolean NOT NULL,
    file_path text NOT NULL,
    source_camera_id integer NOT NULL
);
CREATE SEQUENCE public.video_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.video_id_seq OWNED BY public.video.id;
ALTER TABLE ONLY public.camera ALTER COLUMN id SET DEFAULT nextval('public.camera_id_seq'::regclass);
ALTER TABLE ONLY public.event ALTER COLUMN id SET DEFAULT nextval('public.event_id_seq'::regclass);
ALTER TABLE ONLY public.image ALTER COLUMN id SET DEFAULT nextval('public.image_id_seq'::regclass);
ALTER TABLE ONLY public.video ALTER COLUMN id SET DEFAULT nextval('public.video_id_seq'::regclass);
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
ALTER TABLE ONLY public.video
    ADD CONSTRAINT video_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.video
    ADD CONSTRAINT video_uuid_key UNIQUE (uuid);
ALTER TABLE ONLY public.event
    ADD CONSTRAINT event_high_quality_image_id_fkey FOREIGN KEY (high_quality_image_id) REFERENCES public.image(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.event
    ADD CONSTRAINT event_high_quality_video_id_fkey FOREIGN KEY (high_quality_video_id) REFERENCES public.video(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.event
    ADD CONSTRAINT event_low_quality_image_id_fkey FOREIGN KEY (low_quality_image_id) REFERENCES public.image(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.event
    ADD CONSTRAINT event_low_quality_video_id_fkey FOREIGN KEY (low_quality_video_id) REFERENCES public.video(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.event
    ADD CONSTRAINT event_source_camera_fkey FOREIGN KEY (source_camera_id) REFERENCES public.camera(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.image
    ADD CONSTRAINT image_source_camera_id_fkey FOREIGN KEY (source_camera_id) REFERENCES public.camera(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.video
    ADD CONSTRAINT video_source_camera_id_fkey FOREIGN KEY (source_camera_id) REFERENCES public.camera(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
