package model

import (
	"github.com/google/uuid"
	"github.com/relvacode/iso8601"

	"github.com/initialed85/cameranator/pkg/utils"
)

type Image struct {
	ID             int64        `json:"id,omitempty"`
	UUID           uuid.UUID    `json:"uuid,omitempty"`
	Timestamp      iso8601.Time `json:"timestamp,omitempty"`
	Size           float64      `json:"size,omitempty"`
	FilePath       string       `json:"file_path,omitempty"`
	IsHighQuality  bool         `json:"is_high_quality"`
	SourceCameraID int64        `json:"source_camera_id,omitempty"`
	SourceCamera   Camera       `json:"source_camera,omitempty"`
}

func NewImage(
	Timestamp iso8601.Time,
	Size float64,
	FilePath string,
	IsHighQuality bool,
	SourceCamera Camera,
) Image {
	randomUUID := utils.GetUUID()

	return Image{
		UUID:          randomUUID,
		Timestamp:     Timestamp,
		Size:          Size,
		FilePath:      FilePath,
		IsHighQuality: IsHighQuality,
		SourceCamera:  SourceCamera,
	}
}

func NewImageWithID(
	Timestamp iso8601.Time,
	Size float64,
	FilePath string,
	IsHighQuality bool,
	SourceCameraID int64,
) Image {
	randomUUID := utils.GetUUID()

	return Image{
		UUID:           randomUUID,
		Timestamp:      Timestamp,
		Size:           Size,
		FilePath:       FilePath,
		IsHighQuality:  IsHighQuality,
		SourceCameraID: SourceCameraID,
	}
}
