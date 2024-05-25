package event_receiver

import (
	"encoding/json"
	"testing"
	"time"

	"net"

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
	require.NoError(t, err)
	err = eventReceiver.Open()
	defer eventReceiver.Close()
	require.NoError(t, err)
	time.Sleep(time.Millisecond * 100)

	addr, _ := net.ResolveUDPAddr("udp4", "localhost:6291")

	sender := network.NewSender(addr)
	err = sender.Open()
	defer sender.Close()
	require.NoError(t, err)
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
	require.NoError(t, err)

	err = sender.Send(testEventJSON)
	require.NoError(t, err)

	time.Sleep(time.Millisecond * 1000)

	assert.Equal(
		t,
		testEvent,
		events[len(events)-1],
	)
}
