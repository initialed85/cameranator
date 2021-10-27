CREATE INDEX ON event (start_timestamp);
CREATE INDEX ON event (is_segment, start_timestamp, end_timestamp, source_camera_id);
