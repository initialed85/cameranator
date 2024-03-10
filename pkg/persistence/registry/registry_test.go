package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/initialed85/cameranator/pkg/persistence/model"
)

func TestRegistry_GetModelAndClient(t *testing.T) {
	r := NewRegistry()

	err := r.Register(testGetModel())
	if err != nil {
		require.NoError(t, err)
	}

	modelAndClient, err := r.GetModelAndClient("camera", testGetClient())
	if err != nil {
		require.NoError(t, err)
	}

	items := make([]model.Camera, 0)
	err = modelAndClient.GetAll(&items)
	if err != nil {
		require.NoError(t, err)
	}

	assert.Equal(
		t,
		[]model.Camera{
			{ID: 1, Name: "Driveway", StreamURL: "rtsp://192.168.137.31:554/Streaming/Channels/101/"},
			{ID: 2, Name: "FrontDoor", StreamURL: "rtsp://192.168.137.32:554/Streaming/Channels/101/"},
			{ID: 3, Name: "SideGate", StreamURL: "rtsp://192.168.137.33:554/Streaming/Channels/101/"},
		},
		items,
	)
}
