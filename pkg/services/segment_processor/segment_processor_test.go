package segment_processor

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/initialed85/glue/pkg/network"
	"github.com/relvacode/iso8601"

	"github.com/initialed85/cameranator/pkg/media/converter"
	"github.com/initialed85/cameranator/pkg/persistence/model"
	"github.com/initialed85/cameranator/pkg/segments/segment_generator"
	"github.com/initialed85/cameranator/pkg/utils"
)

func TestSegmentProcessor(t *testing.T) {
	converter.DisableNvidia()

	m, err := NewSegmentProcessor(
		6291,
		"http://localhost:8082/v1/graphql",
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

	rawVideoStartTimestamp := "2020-12-25T08:45:04"
	videoStartTimestamp := utils.GetISO8601Time(rawVideoStartTimestamp)
	videoEndTimestamp := iso8601.Time{Time: videoStartTimestamp.Add(time.Minute * 5)}

	testEvent := segment_generator.Event{
		CameraName:          "Driveway",
		VideoPath:           "../../../test_data/segments/Segment_2020-12-25T08:45:04_Driveway.mp4",
		ImagePath:           "../../../test_data/segments/Segment_2020-12-25T08:45:04_Driveway.jpg",
		VideoStartTimestamp: videoStartTimestamp,
		VideoEndTimestamp:   videoEndTimestamp,
		ImageTimestamp:      videoStartTimestamp,
	}

	testEventJSON, err := json.Marshal(testEvent)
	if err != nil {
		log.Fatal(err)
	}

	err = sender.Send(testEventJSON)
	if err != nil {
		log.Fatal(err)
	}

	timeout := time.Now().Add(time.Second * 30)
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
