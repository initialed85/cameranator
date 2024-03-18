package model

import (
	"github.com/relvacode/iso8601"
)

type Object struct {
	ID             int64        `json:"id,omitempty"`
	StartTimestamp iso8601.Time `json:"start_timestamp,omitempty"`
	EndTimestamp   iso8601.Time `json:"end_timestamp,omitempty"`
	ClassID        int64        `json:"class_id,omitempty"`
	ClassName      string       `json:"class_name,omitempty"`
	CameraID       int64        `json:"camera_id,omitempty"`
	Camera         Camera       `json:"camera,omitempty"`
	EventID        int64        `json:"event_id,omitempty"`
	Event          Event        `json:"event,omitempty"`
}

func NewObject(
	startTimestamp iso8601.Time,
	endTimestamp iso8601.Time,
	classID int64,
	className string,
	camera Camera,
	event Event,
) Object {
	return Object{
		StartTimestamp: startTimestamp,
		EndTimestamp:   endTimestamp,
		ClassID:        classID,
		ClassName:      className,
		Camera:         camera,
		Event:          event,
	}
}

func NewObjectWithIDs(
	startTimestamp iso8601.Time,
	endTimestamp iso8601.Time,
	classID int64,
	className string,
	cameraID int64,
	eventID int64,
) Object {
	return Object{
		StartTimestamp: startTimestamp,
		EndTimestamp:   endTimestamp,
		ClassID:        classID,
		ClassName:      className,
		CameraID:       cameraID,
		EventID:        eventID,
	}
}
