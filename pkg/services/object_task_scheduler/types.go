package object_task_scheduler

import (
	"github.com/google/uuid"
	"github.com/relvacode/iso8601"
)

type PartialVideo struct {
	FilePath string `json:"file_path,omitempty"`
}

type PartialEvent struct {
	UUID             uuid.UUID    `json:"uuid,omitempty"`
	StartTimestamp   iso8601.Time `json:"start_timestamp,omitempty"`
	EndTimestamp     iso8601.Time `json:"end_timestamp,omitempty"`
	HighQualityVideo PartialVideo `json:"high_quality_video,omitempty"`
}
