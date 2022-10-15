INSERT INTO public.camera
    (
        id, uuid, name, stream_url, external_id
    )
VALUES
    (
        1, '3830e9a5-673d-4e7f-ae9b-afa9aeb439ab', 'Driveway', 'rtsp://192.168.137.31:554/Streaming/Channels/101/', '101'
    );
INSERT INTO public.camera
    (
        id, uuid, name, stream_url, external_id
    )
VALUES
    (
        2, 'cd056389-b0b0-4978-9167-68c93e59f53d', 'FrontDoor', 'rtsp://192.168.137.32:554/Streaming/Channels/101/', '102'
    );
INSERT INTO public.camera
    (
        id, uuid, name, stream_url, external_id
    )
VALUES
    (
        3, 'ba9a4013-9c25-4bd0-931f-6eaf61e7369f', 'SideGate', 'rtsp://192.168.137.33:554/Streaming/Channels/101/', '103'
    );
SELECT pg_catalog.setval('public.camera_id_seq', 3, TRUE);
