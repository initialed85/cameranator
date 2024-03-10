package event_receiver

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/initialed85/glue/pkg/network"
	"github.com/relvacode/iso8601"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/initialed85/cameranator/pkg/segments/segment_generator"
	"github.com/initialed85/cameranator/pkg/utils"
)

func TestNewEventReceiver_All(t *testing.T) {
	events := make([]segment_generator.Event, 0)

	eventReceiver, err := NewEventReceiver(
		6291,
		func(event segment_generator.Event) {
			events = append(events, event)
		},
	)
	if err != nil {
		require.NoError(t, err)
	}
	err = eventReceiver.Open()
	defer eventReceiver.Close()
	if err != nil {
		require.NoError(t, err)
	}
	time.Sleep(time.Millisecond * 100)

	sender := network.NewSender("localhost:6291")
	err = sender.Open()
	defer sender.Close()
	if err != nil {
		require.NoError(t, err)
	}
	time.Sleep(time.Millisecond * 100)

	rawVideoStartTimestamp := "2020-12-25T08:45:04"
	videoStartTimestamp := utils.GetISO8601Time(rawVideoStartTimestamp)
	videoEndTimestamp := iso8601.Time{Time: videoStartTimestamp.Add(time.Minute * 5)}

	testEvent := segment_generator.Event{
		CameraName:          "Driveway",
		VideoPath:           "./test_data/segments/Segment_2020-12-25T08:45:04_Driveway.mp4",
		ImagePath:           "./test_data/segments/Segment_2020-12-25T08:45:04_Driveway.jpg",
		VideoStartTimestamp: videoStartTimestamp,
		VideoEndTimestamp:   videoEndTimestamp,
		ImageTimestamp:      videoStartTimestamp,
	}

	testEventJSON, err := json.Marshal(testEvent)
	if err != nil {
		require.NoError(t, err)
	}

	err = sender.Send(testEventJSON)
	if err != nil {
		require.NoError(t, err)
	}

	time.Sleep(time.Millisecond * 1000)

	assert.Equal(
		t,
		[]segment_generator.Event{testEvent},
		events,
	)
}
