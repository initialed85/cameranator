package registry

import (
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/initialed85/cameranator/pkg/common"
	"github.com/initialed85/cameranator/pkg/persistence/graphql"
	"github.com/initialed85/cameranator/pkg/persistence/model"
)

func testGetClient() *graphql.Client {
	return graphql.NewClient("http://localhost:8079/v1/graphql", time.Second*30)
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
		log.Fatal(err)
	}

	for _, camera := range cameras {
		camera.UUID = uuid.UUID{}
	}

	assert.Equal(
		t,
		[]model.Camera{
			{
				ID:        1,
				UUID:      uuid.UUID{0x38, 0x30, 0xe9, 0xa5, 0x67, 0x3d, 0x4e, 0x7f, 0xae, 0x9b, 0xaf, 0xa9, 0xae, 0xb4, 0x39, 0xab},
				Name:      "Driveway",
				StreamURL: "rtsp://192.168.137.31:554/Streaming/Channels/101/",
			},
			{
				ID:        2,
				UUID:      uuid.UUID{0xcd, 0x5, 0x63, 0x89, 0xb0, 0xb0, 0x49, 0x78, 0x91, 0x67, 0x68, 0xc9, 0x3e, 0x59, 0xf5, 0x3d},
				Name:      "FrontDoor",
				StreamURL: "rtsp://192.168.137.32:554/Streaming/Channels/101/",
			},
			{
				ID:        3,
				UUID:      uuid.UUID{0xba, 0x9a, 0x40, 0x13, 0x9c, 0x25, 0x4b, 0xd0, 0x93, 0x1f, 0x6e, 0xaf, 0x61, 0xe7, 0x36, 0x9f},
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

	err := m.GetOne(&cameras, "uuid", "3830e9a5-673d-4e7f-ae9b-afa9aeb439ab")
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(
		t,
		[]model.Camera{
			{
				ID:        1,
				UUID:      uuid.UUID{0x38, 0x30, 0xe9, 0xa5, 0x67, 0x3d, 0x4e, 0x7f, 0xae, 0x9b, 0xaf, 0xa9, 0xae, 0xb4, 0x39, 0xab},
				Name:      "Driveway",
				StreamURL: "rtsp://192.168.137.31:554/Streaming/Channels/101/",
			},
		},
		cameras,
	)
}

func TestModel_AddAndRemove(t *testing.T) {
	m := NewModelAndClient(testGetModel(), testGetClient())

	cameraUUID := utils.GetUUID()

	camera := model.Camera{
		UUID:      cameraUUID,
		Name:      "TestCamera",
		StreamURL: "rtsp://192.168.137.34:554/Streaming/Channels/101/",
	}

	var err error
	cameras := make([]model.Camera, 0)

	err = m.GetAll(&cameras)
	if err != nil {
		log.Fatal(err)
	}
	assert.Len(t, cameras, 3)

	err = m.Add(&camera)
	if err != nil {
		log.Fatal()
	}

	cameras = make([]model.Camera, 0)
	err = m.GetAll(&cameras)
	if err != nil {
		log.Fatal(err)
	}
	assert.Len(t, cameras, 4)

	log.Printf("%#+v", camera)

	err = m.Remove(&camera)
	if err != nil {
		log.Fatal()
	}

	cameras = make([]model.Camera, 0)
	err = m.GetAll(&cameras)
	if err != nil {
		log.Fatal(err)
	}
	assert.Len(t, cameras, 3)
}
