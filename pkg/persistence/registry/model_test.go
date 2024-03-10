package registry

import (
	"log"
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

func TestModel_GetAll(t *testing.T) {
	m := NewModelAndClient(testGetModel(), testGetClient())

	cameras := make([]model.Camera, 0)

	err := m.GetAll(&cameras)
	if err != nil {
		require.NoError(t, err)
	}

	assert.Equal(
		t,
		[]model.Camera{
			{
				ID:        1,
				Name:      "Driveway",
				StreamURL: "rtsp://192.168.137.31:554/Streaming/Channels/101/",
			},
			{
				ID:        2,
				Name:      "FrontDoor",
				StreamURL: "rtsp://192.168.137.32:554/Streaming/Channels/101/",
			},
			{
				ID:        3,
				Name:      "SideGate",
				StreamURL: "rtsp://192.168.137.33:554/Streaming/Channels/101/",
			},
		},
		cameras,
	)
}

func TestModel_GetOne(t *testing.T) {
	m := NewModelAndClient(testGetModel(), testGetClient())

	cameras := make([]model.Camera, 0)

	err := m.GetOne(&cameras, "id", 1)
	if err != nil {
		require.NoError(t, err)
	}

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

func TestModel_AddAndRemove(t *testing.T) {
	m := NewModelAndClient(testGetModel(), testGetClient())

	camera := model.Camera{
		Name:      "TestCamera",
		StreamURL: "rtsp://192.168.137.34:554/Streaming/Channels/101/",
	}

	var err error
	cameras := make([]model.Camera, 0)

	err = m.GetAll(&cameras)
	if err != nil {
		require.NoError(t, err)
	}
	assert.Len(t, cameras, 3)

	defer func() {
		_ = m.Remove(&camera, &cameras)
	}()

	cameras = make([]model.Camera, 0)
	err = m.Add(&camera, &cameras)
	if err != nil {
		require.NoError(t, err)
	}
	log.Printf("%#+v", cameras)

	cameras = make([]model.Camera, 0)
	err = m.GetAll(&cameras)
	if err != nil {
		require.NoError(t, err)
	}
	assert.Len(t, cameras, 4)

	cameras = make([]model.Camera, 0)
	err = m.Remove(&camera, &cameras)
	if err != nil {
		require.NoError(t, err)
	}
	log.Printf("%#+v", cameras)

	cameras = make([]model.Camera, 0)
	err = m.GetAll(&cameras)
	if err != nil {
		require.NoError(t, err)
	}
	assert.Len(t, cameras, 3)
}
