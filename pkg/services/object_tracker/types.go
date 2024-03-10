package object_tracker

import (
	"github.com/relvacode/iso8601"
)

type FFProbeOutput struct {
	Streams []struct {
		Width        int    `json:"width"`
		Height       int    `json:"height"`
		RawFrameRate string `json:"r_frame_rate"`
		RawDuration  string `json:"duration"`
	} `json:"streams"`
}

type PartialVideo struct {
	FilePath       string       `json:"file_path,omitempty"`
	CameraID       int64        `json:"camera_id,omitempty"`
	StartTimestamp iso8601.Time `json:"start_timestamp,omitempty"`
	EndTimestamp   iso8601.Time `json:"end_timestamp,omitempty"`
}

type Point struct {
	X float64
	Y float64
}

type PartialDetection struct {
	Timestamp      iso8601.Time `json:"timestamp,omitempty"`
	RawCentroid    string       `json:"centroid,omitempty"`
	RawBoundingBox string       `json:"bounding_box,omitempty"`
	ClassID        int64        `json:"class_id,omitempty"`
	ClassName      string       `json:"class_name,omitempty"`
	Score          float64      `json:"score,omitempty"`
	Centroid       Point
	BoundingBox    []Point
}

type PartialEvent struct {
	ID               int64              `json:"id,omitempty"`
	OriginalVideo    PartialVideo       `json:"original_video,omitempty"`
	PartialDetection []PartialDetection `json:"detections,omitempty"`
}
