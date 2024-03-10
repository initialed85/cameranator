package model

import (
	"github.com/relvacode/iso8601"
)

type Video struct {
	ID             int64        `json:"id,omitempty"`
	StartTimestamp iso8601.Time `json:"start_timestamp,omitempty"`
	EndTimestamp   iso8601.Time `json:"end_timestamp,omitempty"`
	Size           float64      `json:"size,omitempty"`
	FilePath       string       `json:"file_path,omitempty"`
	CameraID       int64        `json:"camera_id,omitempty"`
	Camera         Camera       `json:"camera,omitempty"`
}

func NewVideo(
	startTimestamp iso8601.Time,
	endTimestamp iso8601.Time,
	size float64,
	filePath string,
	camera Camera,
) Video {
	return Video{
		StartTimestamp: startTimestamp,
		EndTimestamp:   endTimestamp,
		Size:           size,
		FilePath:       filePath,
		Camera:         camera,
	}
}

func NewVideoWithID(
	startTimestamp iso8601.Time,
	endTimestamp iso8601.Time,
	size float64,
	filePath string,
	cameraID int64,
) Video {
	return Video{
		StartTimestamp: startTimestamp,
		EndTimestamp:   endTimestamp,
		Size:           size,
		FilePath:       filePath,
		CameraID:       cameraID,
	}
}
