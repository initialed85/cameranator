package event_pruner

import (
	"log"
	"os"
	"time"

	"github.com/initialed85/glue/pkg/worker"

	"github.com/initialed85/cameranator/pkg/persistence/application"
	"github.com/initialed85/cameranator/pkg/persistence/model"
)

type EventPruner struct {
	scheduledWorker *worker.ScheduledWorker
	application     *application.Application
}

func NewEventPruner(
	url string,
	timeout time.Duration,
	interval time.Duration,
) (*EventPruner, error) {
	var err error

	e := EventPruner{}

	e.scheduledWorker = worker.NewScheduledWorker(
		func() {},
		e.work,
		func() {},
		interval,
	)

	e.application, err = application.NewApplication(url, timeout)

	return &e, err
}

func (e *EventPruner) work() {
	eventModel, err := e.application.GetModelAndClient("event")
	if err != nil {
		log.Printf("warning: %v", err)
		return
	}

	videoModel, err := e.application.GetModelAndClient("video")
	if err != nil {
		log.Printf("warning: %v", err)
		return
	}

	imageModel, err := e.application.GetModelAndClient("image")
	if err != nil {
		log.Printf("warning: %v", err)
		return
	}

	objectModel, err := e.application.GetModelAndClient("object")
	if err != nil {
		log.Printf("warning: %v", err)
		return
	}

	events := make([]model.Event, 0)

	err = eventModel.GetAll(&events)
	if err != nil {
		log.Printf("warning: %v", err)
		return
	}

	for _, event := range events {
		paths := make([]string, 0)
		paths = append(paths, event.OriginalVideo.FilePath)
		paths = append(paths, event.ThumbnailImage.FilePath)

		remove := true
		for _, path := range paths {
			_, err := os.Stat(path)

			if err == nil {
				remove = false
				break
			}
		}

		if !remove {
			continue
		}

		//
		// event
		//

		log.Printf("attempting to delete %#+v", event)
		err = eventModel.Remove(event, &[]model.Event{})
		if err != nil {
			log.Printf("warning: %v", err)
			continue
		}

		//
		// object
		//

		objects := make([]model.Object, 0)
		err = objectModel.GetMany(&objects, "event_id", event.ID)
		if err != nil {
			log.Printf("warning: %v", err)
			continue
		}

		for _, object := range objects {
			err = objectModel.Remove(object, &[]model.Object{})
			if err != nil {
				log.Printf("warning: %v", err)
				continue
			}
		}

		//
		// video
		//

		videos := make([]model.Video, 0)
		err = videoModel.GetOne(&videos, "id", event.OriginalVideoID)
		if err != nil {
			log.Printf("warning: %v", err)
			continue
		}
		video := videos[0]

		log.Printf("attempting to delete %#+v", video)
		err = videoModel.Remove(video, &[]model.Video{})
		if err != nil {
			log.Printf("warning: %v", err)
			continue
		}

		//
		// image
		//

		images := make([]model.Image, 0)
		err = imageModel.GetOne(&images, "id", event.ThumbnailImageID)
		if err != nil {
			log.Printf("warning: %v", err)
			continue
		}
		image := images[0]

		log.Printf("attempting to delete %#+v", image)
		err = imageModel.Remove(image, &[]model.Image{})
		if err != nil {
			log.Printf("warning: %v", err)
			continue
		}
	}
}

func (e *EventPruner) RunOnce() {
	e.work()
}

func (e *EventPruner) Start() {
	e.scheduledWorker.Start()
}

func (e *EventPruner) Stop() {
	e.scheduledWorker.Stop()
}
