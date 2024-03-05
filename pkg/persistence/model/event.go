package model

import (
	"github.com/google/uuid"
	"github.com/relvacode/iso8601"

	"github.com/initialed85/cameranator/pkg/utils"
)

type Event struct {
	ID                    int64        `json:"id,omitempty"`
	UUID                  uuid.UUID    `json:"uuid,omitempty"`
	StartTimestamp        iso8601.Time `json:"start_timestamp,omitempty"`
	EndTimestamp          iso8601.Time `json:"end_timestamp,omitempty"`
	IsSegment             bool         `json:"is_segment"`
	HighQualityVideoID    int64        `json:"high_quality_video_id,omitempty"`
	HighQualityVideo      Video        `json:"high_quality_video,omitempty"`
	HighQualityImageID    int64        `json:"high_quality_image_id,omitempty"`
	HighQualityImage      Image        `json:"high_quality_image,omitempty"`
	LowQualityVideoID     int64        `json:"low_quality_video_id,omitempty"`
	LowQualityVideo       Video        `json:"low_quality_video,omitempty"`
	LowQualityImageID     int64        `json:"low_quality_image_id,omitempty"`
	LowQualityImage       Image        `json:"low_quality_image,omitempty"`
	SourceCameraID        int64        `json:"source_camera_id,omitempty"`
	SourceCamera          Camera       `json:"source_camera,omitempty"`
	NeedsObjectProcessing bool         `json:"needs_object_processing"`
}

func NewEvent(
	StartTimestamp iso8601.Time,
	EndTimestamp iso8601.Time,
	IsSegment bool,
	HighQualityVideo Video,
	HighQualityImage Image,
	LowQualityVideo Video,
	LowQualityImage Image,
	SourceCamera Camera,
) Event {
	randomUUID := utils.GetUUID()

	return Event{
		UUID:             randomUUID,
		StartTimestamp:   StartTimestamp,
		EndTimestamp:     EndTimestamp,
		IsSegment:        IsSegment,
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
	IsSegment bool,
	HighQualityVideoID int64,
	HighQualityImageID int64,
	LowQualityVideoID int64,
	LowQualityImageID int64,
	SourceCameraID int64,
) Event {
	randomUUID := utils.GetUUID()

	return Event{
		UUID:               randomUUID,
		StartTimestamp:     StartTimestamp,
		EndTimestamp:       EndTimestamp,
		IsSegment:          IsSegment,
		HighQualityVideoID: HighQualityVideoID,
		HighQualityImageID: HighQualityImageID,
		LowQualityVideoID:  LowQualityVideoID,
		LowQualityImageID:  LowQualityImageID,
		SourceCameraID:     SourceCameraID,
	}
}
