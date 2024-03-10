package segment_generators

import (
	"encoding/json"
	"log"
	"net"
	"testing"
	"time"

	"github.com/initialed85/glue/pkg/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/initialed85/cameranator/pkg/segments/segment_generator"
	"github.com/initialed85/cameranator/pkg/test_utils"
)

func TestNewSegmentGenerators(t *testing.T) {
	dir, err := test_utils.GetTempDir()
	if err != nil {
		require.NoError(t, err)
	}

	interfaceName, err := network.GetDefaultInterfaceName()
	if err != nil {
		require.NoError(t, err)
	}

	receiver := network.NewReceiver("0.0.0.0:6291", interfaceName)

	events := make([]segment_generator.Event, 0)

	err = receiver.RegisterCallback(func(addr *net.UDPAddr, data []byte) {
		event := segment_generator.Event{}

		err := json.Unmarshal(data, &event)
		if err != nil {
			require.NoError(t, err)
		}

		events = append(events, event)
	})
	if err != nil {
		require.NoError(t, err)
	}

	err = receiver.Open()
	if err != nil {
		require.NoError(t, err)
	}
	defer receiver.Close()

	time.Sleep(time.Millisecond * 100)

	segmentGenerators := NewSegmentGenerators(
		[]segment_generator.Feed{
			{
				NetCamURL:       "rtsp://host.docker.internal:8554/Streaming/Channels/101",
				DestinationPath: dir,
				CameraName:      "Testing",
				Duration:        5,
			},
		},
		"localhost",
		6291,
	)

	err = segmentGenerators.Start()
	defer segmentGenerators.Stop()
	if err != nil {
		require.NoError(t, err)
	}

	assert.True(t, test_utils.IsLive("localhost", 8080))

	time.Sleep(time.Second * 15)

	segmentGenerators.Stop()

	time.Sleep(time.Second * 15)

	assert.False(t, test_utils.IsLive("localhost", 8080))

	assert.Greater(t, len(events), 0)

	log.Printf("%#+v", events)
}
