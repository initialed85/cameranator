package object_task_scheduler

import "github.com/relvacode/iso8601"

type PartialVideo struct {
	FilePath       string       `json:"file_path,omitempty"`
	SourceCameraID int64        `json:"source_camera_id,omitempty"`
	StartTimestamp iso8601.Time `json:"start_timestamp,omitempty"`
	EndTimestamp   iso8601.Time `json:"end_timestamp,omitempty"`
}

type PartialEvent struct {
	ID               int64        `json:"id,omitempty"`
	HighQualityVideo PartialVideo `json:"high_quality_video,omitempty"`
}
