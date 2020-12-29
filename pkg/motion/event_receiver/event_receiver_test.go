package event_receiver

import (
	"log"
	"testing"
	"time"

	"github.com/initialed85/glue/pkg/network"
	"github.com/relvacode/iso8601"
	"github.com/stretchr/testify/assert"
)

var testMessages = []string{
	`[2020-12-27T10:25:10,034772700+08:00, on_event_start, 2020-12-27T10:25:09, 104, Testing, 01, null]`,
	`[2020-12-27T10:25:10,037680700+08:00, on_movie_start, 2020-12-27T10:25:09, 104, Testing, 01, ./test_data/events/Event_2020-12-27T10:25:05__104__Testing__01.mp4]`,
	`[2020-12-27T10:25:10,141712800+08:00, on_picture_save, 2020-12-27T10:25:09, 104, Testing, 01, ./test_data/events/Event_2020-12-27T10:25:09__104__Testing__01.jpg]`,
	`[2020-12-27T10:25:21,179423900+08:00, on_event_end, 2020-12-27T10:25:21, 104, Testing, 01, null]`,
	`[2020-12-27T10:25:21,188467400+08:00, on_movie_end, 2020-12-27T10:25:21, 104, Testing, 01, ./test_data/events/Event_2020-12-27T10:25:05__104__Testing__01.mp4]`,
	`[2020-12-27T10:25:27,458801900+08:00, on_event_start, 2020-12-27T10:25:27, 104, Testing, 02, null]`,
	`[2020-12-27T10:25:27,470455700+08:00, on_movie_start, 2020-12-27T10:25:27, 104, Testing, 02, ./test_data/events/Event_2020-12-27T10:25:23__104__Testing__02.mp4]`,
	`[2020-12-27T10:25:27,529506200+08:00, on_picture_save, 2020-12-27T10:25:27, 104, Testing, 02, ./test_data/events/Event_2020-12-27T10:25:27__104__Testing__02.jpg]`,
	`[2020-12-27T10:25:39,160444400+08:00, on_event_end, 2020-12-27T10:25:39, 104, Testing, 02, null]`,
	`[2020-12-27T10:25:39,184958000+08:00, on_movie_end, 2020-12-27T10:25:39, 104, Testing, 02, ./test_data/events/Event_2020-12-27T10:25:23__104__Testing__02.mp4]`,
	`[2020-12-27T10:26:14,629448600+08:00, on_event_start, 2020-12-27T10:26:14, 104, Testing, 03, null]`,
	`[2020-12-27T10:26:14,673063600+08:00, on_movie_start, 2020-12-27T10:26:14, 104, Testing, 03, ./test_data/events/Event_2020-12-27T10:26:10__104__Testing__03.mp4]`,
	`[2020-12-27T10:26:14,729056300+08:00, on_picture_save, 2020-12-27T10:26:14, 104, Testing, 03, ./test_data/events/Event_2020-12-27T10:26:14__104__Testing__03.jpg]`,
	`[2020-12-27T10:26:38,170122700+08:00, on_event_end, 2020-12-27T10:26:38, 104, Testing, 03, null]`,
	`[2020-12-27T10:26:38,192276800+08:00, on_movie_end, 2020-12-27T10:26:38, 104, Testing, 03, ./test_data/events/Event_2020-12-27T10:26:10__104__Testing__03.mp4]`,
	`[2020-12-27T10:29:27,717699100+08:00, on_event_start, 2020-12-27T10:29:27, 104, Testing, 04, null]`,
	`[2020-12-27T10:29:27,756390800+08:00, on_movie_start, 2020-12-27T10:29:27, 104, Testing, 04, ./test_data/events/Event_2020-12-27T10:29:23__104__Testing__04.mp4]`,
	`[2020-12-27T10:29:27,832972500+08:00, on_picture_save, 2020-12-27T10:29:27, 104, Testing, 04, ./test_data/events/Event_2020-12-27T10:29:27__104__Testing__04.jpg]`,
	`[2020-12-27T10:29:47,179383300+08:00, on_event_end, 2020-12-27T10:29:47, 104, Testing, 04, null]`,
	`[2020-12-27T10:29:47,189592900+08:00, on_movie_end, 2020-12-27T10:29:47, 104, Testing, 04, ./test_data/events/Event_2020-12-27T10:29:23__104__Testing__04.mp4]`,
	`[2020-12-27T10:30:01,136402900+08:00, on_event_end, 2020-12-27T10:30:00, 104, Testing, 05, null]`,
}

func TestNewEventReceiver_All(t *testing.T) {
	events := make([]Event, 0)

	eventReceiver, err := NewEventReceiver(
		6291,
		func(event Event) {
			events = append(events, event)
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	err = eventReceiver.Open()
	defer eventReceiver.Close()
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

		time.Sleep(time.Millisecond * 10)
	}

	time.Sleep(time.Millisecond * 1000)

	for i, event := range events {
		assert.True(t, event.StartTimestamp.After(time.Time{}))
		assert.True(t, event.VideoStartTimestamp.After(time.Time{}))
		assert.True(t, event.ImageTimestamp.After(time.Time{}))
		assert.True(t, event.VideoEndTimestamp.After(event.VideoStartTimestamp.Time))
		assert.True(t, event.EndTimestamp.After(event.StartTimestamp.Time))

		events[i].StartTimestamp = iso8601.Time{}
		events[i].VideoStartTimestamp = iso8601.Time{}
		events[i].ImageTimestamp = iso8601.Time{}
		events[i].VideoEndTimestamp = iso8601.Time{}
		events[i].EndTimestamp = iso8601.Time{}
	}

	assert.Equal(
		t,
		[]Event{
			{
				CameraID:            104,
				CameraName:          "Testing",
				EventNumber:         1,
				VideoPath:           "./test_data/events/Event_2020-12-27T10:25:05__104__Testing__01.mp4",
				ImagePath:           "./test_data/events/Event_2020-12-27T10:25:09__104__Testing__01.jpg",
				StartTimestamp:      iso8601.Time{},
				VideoStartTimestamp: iso8601.Time{},
				ImageTimestamp:      iso8601.Time{},
				VideoEndTimestamp:   iso8601.Time{},
				EndTimestamp:        iso8601.Time{},
				Complete:            true,
			},
			{
				CameraID:            104,
				CameraName:          "Testing",
				EventNumber:         2,
				VideoPath:           "./test_data/events/Event_2020-12-27T10:25:23__104__Testing__02.mp4",
				ImagePath:           "./test_data/events/Event_2020-12-27T10:25:27__104__Testing__02.jpg",
				StartTimestamp:      iso8601.Time{},
				VideoStartTimestamp: iso8601.Time{},
				ImageTimestamp:      iso8601.Time{},
				VideoEndTimestamp:   iso8601.Time{},
				EndTimestamp:        iso8601.Time{},
				Complete:            true,
			},
			{
				CameraID:            104,
				CameraName:          "Testing",
				EventNumber:         3,
				VideoPath:           "./test_data/events/Event_2020-12-27T10:26:10__104__Testing__03.mp4",
				ImagePath:           "./test_data/events/Event_2020-12-27T10:26:14__104__Testing__03.jpg",
				StartTimestamp:      iso8601.Time{},
				VideoStartTimestamp: iso8601.Time{},
				ImageTimestamp:      iso8601.Time{},
				VideoEndTimestamp:   iso8601.Time{},
				EndTimestamp:        iso8601.Time{},
				Complete:            true,
			},
			{
				CameraID:            104,
				CameraName:          "Testing",
				EventNumber:         4,
				VideoPath:           "./test_data/events/Event_2020-12-27T10:29:23__104__Testing__04.mp4",
				ImagePath:           "./test_data/events/Event_2020-12-27T10:29:27__104__Testing__04.jpg",
				StartTimestamp:      iso8601.Time{},
				VideoStartTimestamp: iso8601.Time{},
				ImageTimestamp:      iso8601.Time{},
				VideoEndTimestamp:   iso8601.Time{},
				EndTimestamp:        iso8601.Time{},
				Complete:            true,
			},
		},
		events,
	)
}
