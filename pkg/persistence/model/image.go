package model

import (
	"github.com/relvacode/iso8601"
)

type Image struct {
	ID        int64        `json:"id,omitempty"`
	Timestamp iso8601.Time `json:"timestamp,omitempty"`
	Size      float64      `json:"size,omitempty"`
	FilePath  string       `json:"file_path,omitempty"`
	CameraID  int64        `json:"camera_id,omitempty"`
	Camera    Camera       `json:"camera,omitempty"`
}

func NewImage(
	timestamp iso8601.Time,
	size float64,
	filePath string,
	camera Camera,
) Image {
	return Image{
		Timestamp: timestamp,
		Size:      size,
		FilePath:  filePath,
		Camera:    camera,
	}
}

func NewImageWithID(
	timestamp iso8601.Time,
	size float64,
	filePath string,
	cameraID int64,
) Image {
	return Image{
		Timestamp: timestamp,
		Size:      size,
		FilePath:  filePath,
		CameraID:  cameraID,
	}
}
