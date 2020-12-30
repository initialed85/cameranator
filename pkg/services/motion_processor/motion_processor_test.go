package motion_processor

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/initialed85/glue/pkg/network"

	"github.com/initialed85/cameranator/pkg/media/converter"
	"github.com/initialed85/cameranator/pkg/persistence/model"
)

var testMessages = []string{
	`[2020-12-27T10:25:10,034772700+08:00, on_event_start, 2020-12-27T10:25:09, 104, Testing, 01, null]`,
	`[2020-12-27T10:25:10,037680700+08:00, on_movie_start, 2020-12-27T10:25:09, 104, Testing, 01, ../../../test_data/events/Event_2020-12-27T10:25:05__104__Testing__01.mp4]`,
	`[2020-12-27T10:25:10,141712800+08:00, on_picture_save, 2020-12-27T10:25:09, 104, Testing, 01, ../../../test_data/events/Event_2020-12-27T10:25:09__104__Testing__01.jpg]`,
	`[2020-12-27T10:25:21,179423900+08:00, on_event_end, 2020-12-27T10:25:21, 104, Testing, 01, null]`,
	`[2020-12-27T10:25:21,188467400+08:00, on_movie_end, 2020-12-27T10:25:21, 104, Testing, 01, ../../../test_data/events/Event_2020-12-27T10:25:05__104__Testing__01.mp4]`,
}

func TestNewMotionProcessor(t *testing.T) {
	converter.DisableNvidia()

	m, err := NewMotionProcessor(
		6291,
		"http://localhost:8079/v1/graphql",
		time.Second*10,
	)
	if err != nil {
		log.Fatal()
	}

	cameraModelAndClient, err := m.application.GetModelAndClient("camera")
	if err != nil {
		log.Fatal(err)
	}

	imageModelAndClient, err := m.application.GetModelAndClient("image")
	if err != nil {
		log.Fatal(err)
	}

	videoModelAndClient, err := m.application.GetModelAndClient("video")
	if err != nil {
		log.Fatal(err)
	}

	eventModelAndClient, err := m.application.GetModelAndClient("event")
	if err != nil {
		log.Fatal(err)
	}

	camera := model.NewCamera(
		"Testing",
		"rtsp://host.docker.internal:8554/Streaming/Channels/101",
	)

	err = cameraModelAndClient.Add(camera, []model.Camera{})
	defer func() {
		events := make([]model.Event, 0)
		err = eventModelAndClient.GetAll(&events)
		if err != nil {
			log.Fatal(err)
		}

		for _, event := range events {
			err = eventModelAndClient.Remove(&event, []model.Event{})
			if err != nil {
				log.Fatal(err)
			}
		}

		videos := make([]model.Video, 0)
		err = videoModelAndClient.GetAll(&videos)
		if err != nil {
			log.Fatal(err)
		}

		for _, video := range videos {
			err = videoModelAndClient.Remove(&video, []model.Video{})
			if err != nil {
				log.Fatal(err)
			}
		}

		images := make([]model.Image, 0)
		err = imageModelAndClient.GetAll(&images)
		if err != nil {
			log.Fatal(err)
		}

		for _, image := range images {
			err = imageModelAndClient.Remove(&image, []model.Image{})
			if err != nil {
				log.Fatal(err)
			}
		}

		err = cameraModelAndClient.Remove(&camera, []model.Camera{})
		if err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		log.Fatal(err)
	}

	err = m.Start()
	defer m.Stop()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Millisecond * 100)

	sender := network.NewSender("localhost:6291")
	err = sender.Open()
	defer sender.Close()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Millisecond * 100)

	for _, message := range testMessages {
		err = sender.Send([]byte(message))
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Millisecond)
	}

	timeout := time.Now().Add(time.Second * 10)
	for time.Now().Before(timeout) {
		events := make([]model.Event, 0)

		err = eventModelAndClient.GetAll(&events)
		if err != nil {
			log.Fatal(err)
		}

		if len(events) > 0 {
			eventsJSON, err := json.MarshalIndent(events, "", "   ")
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("added event %v", string(eventsJSON))

			return
		}
	}

	log.Fatal("timed out")
}
