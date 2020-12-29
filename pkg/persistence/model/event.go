package model

import (
	"github.com/google/uuid"
	"github.com/relvacode/iso8601"

	"github.com/initialed85/cameranator/pkg/common"
)

type Event struct {
	ID                 int64        `json:"id,omitempty"`
	UUID               uuid.UUID    `json:"uuid,omitempty"`
	StartTimestamp     iso8601.Time `json:"start_timestamp,omitempty"`
	EndTimestamp       iso8601.Time `json:"end_timestamp,omitempty"`
	IsProcessed        bool         `json:"is_processed"`
	HighQualityVideoID int64        `json:"high_quality_video_id,omitempty"`
	HighQualityVideo   Video        `json:"high_quality_video,omitempty"`
	HighQualityImageID int64        `json:"high_quality_image_id,omitempty"`
	HighQualityImage   Image        `json:"high_quality_image,omitempty"`
	LowQualityVideoID  int64        `json:"low_quality_video_id,omitempty"`
	LowQualityVideo    Video        `json:"low_quality_video,omitempty"`
	LowQualityImageID  int64        `json:"low_quality_image_id,omitempty"`
	LowQualityImage    Image        `json:"low_quality_image,omitempty"`
	SourceCameraID     int64        `json:"source_camera_id,omitempty"`
	SourceCamera       Camera       `json:"source_camera,omitempty"`
}

func NewEvent(
	StartTimestamp iso8601.Time,
	EndTimestamp iso8601.Time,
	IsProcessed bool,
	HighQualityVideo Video,
	HighQualityImage Image,
	LowQualityVideo Video,
	LowQualityImage Image,
	SourceCamera Camera,
) Event {
	randomUUID := common.GetUUID()

	return Event{
		UUID:             randomUUID,
		StartTimestamp:   StartTimestamp,
		EndTimestamp:     EndTimestamp,
		IsProcessed:      IsProcessed,
		HighQualityVideo: HighQualityVideo,
		HighQualityImage: HighQualityImage,
		LowQualityVideo:  LowQualityVideo,
		LowQualityImage:  LowQualityImage,
		SourceCamera:     SourceCamera,
	}
}

func NewEventWithIDs(
	StartTimestamp iso8601.Time,
	EndTimestamp iso8601.Time,
	IsProcessed bool,
	HighQualityVideoID int64,
	HighQualityImageID int64,
	LowQualityVideoID int64,
	LowQualityImageID int64,
	SourceCameraID int64,
) Event {
	randomUUID := common.GetUUID()

	return Event{
		UUID:               randomUUID,
		StartTimestamp:     StartTimestamp,
		EndTimestamp:       EndTimestamp,
		IsProcessed:        IsProcessed,
		HighQualityVideoID: HighQualityVideoID,
		HighQualityImageID: HighQualityImageID,
		LowQualityVideoID:  LowQualityVideoID,
		LowQualityImageID:  LowQualityImageID,
		SourceCameraID:     SourceCameraID,
	}
}
