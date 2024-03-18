package object_tracker

import (
	"time"

	"github.com/relvacode/iso8601"
)

type FFProbeOutput struct {
	Streams []struct {
		Width       int    `json:"width"`
		Height      int    `json:"height"`
		RawFPS      string `json:"r_frame_rate"`
		RawFrames   string `json:"nb_frames"`
		RawDuration string `json:"duration"`
	} `json:"streams"`
}

type PartialVideo struct {
	ID               int64        `json:"id,omitempty"`
	FilePath         string       `json:"file_path,omitempty"`
	CameraID         int64        `json:"camera_id,omitempty"`
	StartTimestamp   iso8601.Time `json:"start_timestamp,omitempty"`
	EndTimestamp     iso8601.Time `json:"end_timestamp,omitempty"`
	AdjustedFilePath string
	Width            int64
	Height           int64
	FPS              float64
	Frames           int64
	Duration         time.Duration
}

type Point struct {
	X float64
	Y float64
}

type PartialDetection struct {
	ID                                int64        `json:"id,omitempty"`
	Timestamp                         iso8601.Time `json:"timestamp,omitempty"`
	RawCentroid                       string       `json:"centroid,omitempty"`
	RawBoundingBox                    string       `json:"bounding_box,omitempty"`
	ClassID                           int64        `json:"class_id,omitempty"`
	ClassName                         string       `json:"class_name,omitempty"`
	Score                             float64      `json:"score,omitempty"`
	Centroid                          Point
	BoundingBox                       []Point
	Height                            float64
	Width                             float64
	Area                              float64
	AspectRatio                       float64
	Frame                             int64
	SameClassPartialDetectionsByFrame map[int64][]*PartialDetection
	ObjectID                          int64
}

type PartialEvent struct {
	ID                       int64               `json:"id,omitempty"`
	OriginalVideo            PartialVideo        `json:"original_video,omitempty"`
	PartialDetections        []*PartialDetection `json:"detections,omitempty"`
	PartialDetectionsByFrame map[int64][]*PartialDetection
}

type Score struct {
	ID                     int64
	FrameDistanceFactor    float64
	AreaFactor             float64
	AspectRatioFactor      float64
	CentroidDistanceFactor float64
}
