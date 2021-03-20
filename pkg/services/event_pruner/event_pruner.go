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

	events := make([]model.Event, 0)

	err = eventModel.GetAll(&events)
	if err != nil {
		log.Printf("warning: %v", err)
		return
	}

	for _, event := range events {
		paths := make([]string, 0)

		paths = append(paths, event.HighQualityVideo.FilePath)
		paths = append(paths, event.LowQualityVideo.FilePath)
		paths = append(paths, event.HighQualityImage.FilePath)
		paths = append(paths, event.LowQualityImage.FilePath)

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

		log.Printf("attempting to delete %#+v", event)
		err = eventModel.Remove(event, &[]model.Event{})
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
