package model

import (
	"github.com/google/uuid"
	"github.com/relvacode/iso8601"

	"github.com/initialed85/cameranator/pkg/common"
)

type Video struct {
	ID             int64        `json:"id,omitempty"`
	UUID           uuid.UUID    `json:"uuid,omitempty"`
	StartTimestamp iso8601.Time `json:"start_timestamp,omitempty"`
	EndTimestamp   iso8601.Time `json:"end_timestamp,omitempty"`
	Size           float64      `json:"size,omitempty"`
	FilePath       string       `json:"file_path,omitempty"`
	IsHighQuality  bool         `json:"is_high_quality"`
	SourceCameraID int64        `json:"source_camera_id,omitempty"`
	SourceCamera   Camera       `json:"source_camera,omitempty"`
}

func NewVideo(
	StartTimestamp iso8601.Time,
	EndTimestamp iso8601.Time,
	Size float64,
	FilePath string,
	IsHighQuality bool,
	SourceCamera Camera,
) Video {
	randomUUID := utils.GetUUID()

	return Video{
		UUID:           randomUUID,
		StartTimestamp: StartTimestamp,
		EndTimestamp:   EndTimestamp,
		Size:           Size,
		FilePath:       FilePath,
		IsHighQuality:  IsHighQuality,
		SourceCamera:   SourceCamera,
	}
}

func NewVideoWithID(
	StartTimestamp iso8601.Time,
	EndTimestamp iso8601.Time,
	Size float64,
	FilePath string,
	IsHighQuality bool,
	SourceCameraID int64,
) Video {
	randomUUID := utils.GetUUID()

	return Video{
		UUID:           randomUUID,
		StartTimestamp: StartTimestamp,
		EndTimestamp:   EndTimestamp,
		Size:           Size,
		FilePath:       FilePath,
		IsHighQuality:  IsHighQuality,
		SourceCameraID: SourceCameraID,
	}
}
