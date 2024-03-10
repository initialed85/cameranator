package model

import (
	"github.com/relvacode/iso8601"
)

type Event struct {
	ID               int64        `json:"id,omitempty"`
	StartTimestamp   iso8601.Time `json:"start_timestamp,omitempty"`
	EndTimestamp     iso8601.Time `json:"end_timestamp,omitempty"`
	OriginalVideoID  int64        `json:"original_video_id,omitempty"`
	OriginalVideo    Video        `json:"original_video,omitempty"`
	ThumbnailImageID int64        `json:"thumbnail_image_id,omitempty"`
	ThumnailImage    Image        `json:"thumbnail_image,omitempty"`
	SourceCameraID   int64        `json:"source_camera_id,omitempty"`
	SourceCamera     Camera       `json:"source_camera,omitempty"`
	Status           string       `json:"status,omitempty"`
}

func NewEvent(
	startTimestamp iso8601.Time,
	endTimestamp iso8601.Time,
	originalVideo Video,
	thumbnailImage Image,
	sourceCamera Camera,
) Event {
	return Event{
		StartTimestamp: startTimestamp,
		EndTimestamp:   endTimestamp,
		OriginalVideo:  originalVideo,
		ThumnailImage:  thumbnailImage,
		SourceCamera:   sourceCamera,
		Status:         "needs detection",
	}
}

func NewEventWithIDs(
	startTimestamp iso8601.Time,
	endTimestamp iso8601.Time,
	originalVideo int64,
	thumbnailImage int64,
	sourceCameraID int64,
) Event {
	return Event{
		StartTimestamp:   startTimestamp,
		EndTimestamp:     endTimestamp,
		OriginalVideoID:  originalVideo,
		ThumbnailImageID: thumbnailImage,
		SourceCameraID:   sourceCameraID,
		Status:           "needs detection",
	}
}
