package helpers

import (
	"fmt"

	"github.com/relvacode/iso8601"

	"github.com/initialed85/cameranator/pkg/media/metadata"
	"github.com/initialed85/cameranator/pkg/persistence/application"
	"github.com/initialed85/cameranator/pkg/persistence/model"
)

func GetCamera(
	application *application.Application,
	name string,
) (model.Camera, error) {
	cameraModelAndClient, err := application.GetModelAndClient("camera")
	if err != nil {
		return model.Camera{}, err
	}

	cameras := make([]model.Camera, 0)
	err = cameraModelAndClient.GetOne(&cameras, "name", name)
	if err != nil {
		return model.Camera{}, err
	}

	if len(cameras) != 1 {
		return model.Camera{}, fmt.Errorf("failed to find exactly one camera for %#+v", name)
	}

	return cameras[0], nil
}

func AddEvent(
	application *application.Application,
	cameraName string,
	startTimestamp iso8601.Time,
	endTimestamp iso8601.Time,
	highQualityVideoPath string,
	highQualityImagePath string,
	lowQualityVideoPath string,
	lowQualityImagePath string,
) (model.Event, error) {
	camera, err := GetCamera(application, cameraName)
	if err != nil {
		return model.Event{}, err
	}

	highQualityVideoSize, err := metadata.GetFileSize(highQualityVideoPath)
	if err != nil {
		return model.Event{}, err
	}

	highQualityImageSize, err := metadata.GetFileSize(highQualityImagePath)
	if err != nil {
		return model.Event{}, err
	}

	lowQualityVideoSize, err := metadata.GetFileSize(lowQualityVideoPath)
	if err != nil {
		return model.Event{}, err
	}

	lowQualityImageSize, err := metadata.GetFileSize(lowQualityImagePath)
	if err != nil {
		return model.Event{}, err
	}

	highQualityVideo := model.NewVideo(
		startTimestamp,
		endTimestamp,
		highQualityVideoSize,
		highQualityImagePath,
		true,
		camera,
	)

	highQualityImage := model.NewImage(
		startTimestamp,
		highQualityImageSize,
		highQualityImagePath,
		true,
		camera,
	)

	lowQualityVideo := model.NewVideo(
		startTimestamp,
		endTimestamp,
		lowQualityVideoSize,
		lowQualityImagePath,
		true,
		camera,
	)

	lowQualityImage := model.NewImage(
		startTimestamp,
		lowQualityImageSize,
		lowQualityImagePath,
		true,
		camera,
	)

	event := model.NewEvent(
		startTimestamp,
		endTimestamp,
		true,
		highQualityVideo,
		highQualityImage,
		lowQualityVideo,
		lowQualityImage,
		camera,
	)

	eventModelAndClient, err := application.GetModelAndClient("event")
	if err != nil {
		return model.Event{}, err
	}

	events := make([]model.Event, 0)
	err = eventModelAndClient.Add(&event, &events)
	if err != nil {
		return model.Event{}, err
	}

	if len(events) != 1 {
		return model.Event{}, fmt.Errorf("attempt to add Event should have returned exactly 1 Event")
	}

	return events[0], nil
}
