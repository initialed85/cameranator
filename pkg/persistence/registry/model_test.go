package registry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/initialed85/cameranator/pkg/persistence/graphql"
	"github.com/initialed85/cameranator/pkg/persistence/model"
)

func testGetClient() *graphql.Client {
	return graphql.NewClient("http://localhost:8082/v1/graphql", time.Second*30)
}

func testGetModel() *Model {
	m := NewModel("camera", model.Camera{})

	return m
}

func TestModel_AddAndRemove(t *testing.T) {
	m := NewModelAndClient(testGetModel(), testGetClient())

	camera := model.Camera{
		Name:      "TestCamera_TestModel",
		StreamURL: "rtsp://192.168.137.34:554/Streaming/Channels/101/",
	}

	var err error
	cameras := make([]model.Camera, 0)

	err = m.GetAll(&cameras)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(cameras), 0)

	cameras = make([]model.Camera, 0)
	err = m.Add(&camera, &cameras)
	require.NoError(t, err)
	insertedCamera := cameras[0]

	cameras = make([]model.Camera, 0)
	err = m.GetAll(&cameras)
	require.NoError(t, err)
	assert.Condition(t, func() bool {
		for _, thisCamera := range cameras {
			if thisCamera.ID == insertedCamera.ID {
				return true
			}
		}

		return false
	})

	cameras = make([]model.Camera, 0)
	err = m.Remove(&camera, &[]model.Camera{})
	require.NoError(t, err)

	cameras = make([]model.Camera, 0)
	err = m.GetAll(&cameras)
	require.NoError(t, err)
	assert.Condition(t, func() bool {
		for _, thisCamera := range cameras {
			if thisCamera.ID == insertedCamera.ID {
				return false
			}
		}

		return true
	})
}

func TestModel_GetAll(t *testing.T) {
	m := NewModelAndClient(testGetModel(), testGetClient())

	cameras := make([]model.Camera, 0)

	err := m.GetAll(&cameras)
	require.NoError(t, err)

	assert.Condition(t, func() bool {
		for _, camera := range cameras {
			if camera.Name == "Driveway" && camera.StreamURL == "rtsp://192.168.137.31:554/Streaming/Channels/101/" {
				return true
			}
		}
		return false
	})

	assert.Condition(t, func() bool {
		for _, camera := range cameras {
			if camera.Name == "FrontDoor" && camera.StreamURL == "rtsp://192.168.137.32:554/Streaming/Channels/101/" {
				return true
			}
		}
		return false
	})

	assert.Condition(t, func() bool {
		for _, camera := range cameras {
			if camera.Name == "SideGate" && camera.StreamURL == "rtsp://192.168.137.33:554/Streaming/Channels/101/" {
				return true
			}
		}
		return false
	})
}

func TestModel_GetOne(t *testing.T) {
	m := NewModelAndClient(testGetModel(), testGetClient())

	cameras := make([]model.Camera, 0)

	err := m.GetOne(&cameras, "id", 1)
	require.NoError(t, err)

	assert.Equal(
		t,
		[]model.Camera{
			{
				ID:        1,
				Name:      "Driveway",
				StreamURL: "rtsp://192.168.137.31:554/Streaming/Channels/101/",
			},
		},
		cameras,
	)
}
