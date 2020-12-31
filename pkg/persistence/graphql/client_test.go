package graphql

import (
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/initialed85/cameranator/pkg/persistence/model"
)

func testGetClient() *Client {
	return NewClient("http://localhost:8082/v1/graphql", time.Second*30)
}

func testGetManyQuery() string {
	return `
{
  camera (order_by: {id: asc}) {
    uuid
    name
    stream_url
  }
}
	`
}

func testGetOneQuery() string {
	return `
{
  camera(where: {uuid: {_eq: "3830e9a5-673d-4e7f-ae9b-afa9aeb439ab"}}, limit: 1, distinct_on: uuid) {
    uuid
    name
    stream_url
  }
}
`
}

func testInsertOneQuery() string {
	return `
mutation {
  insert_camera_one(object: {uuid: "64dbac5a-29c7-4244-b297-0c540af329f9", name: "TestCamera", stream_url: "rtsp://192.168.137.34:554/Streaming/Channels/101/"}) {
    uuid
    name
    stream_url
  }
}
`
}

func testInsertOneImage() string {
	return `
mutation {
  insert_image_one(object: {
    timestamp: "2020-12-26T01:59:59+00:00",
    size: 65536,
    is_high_quality: true,
    file_path: "path/to/file",
    source_camera_id: 1
  }) {
    uuid
  }
}
`
}

func testInsertOneVideo() string {
	return `
mutation {
  insert_video_one(object: {
    start_timestamp: "2020-12-26T01:59:59+00:00", 
    end_timestamp: "2020-12-26T01:59:59+00:00",
    size: 65536,
    is_high_quality: true,
    file_path: "path/to/file",
    source_camera_id: 1
  }) {
    uuid
  }
}
`
}

func testDeleteQuery() string {
	return `
mutation {
  delete_camera(where: {uuid: {_eq: "64dbac5a-29c7-4244-b297-0c540af329f9"}}) {
    returning {
      uuid
      name
      stream_url
    }
  }
}
`
}

func testGetMultipleManyQuery() string {
	return `
{
  image (order_by: {id: asc}) {
    uuid
    timestamp
    size
    is_high_quality
    file_path
    source_camera {
      uuid
      name
      stream_url
    }
  }
  video (order_by: {id: asc}) {
    uuid
    start_timestamp
    end_timestamp
    size
    is_high_quality
    file_path
    source_camera {
      uuid
      name
      stream_url
    }
  }
}
	`
}

func TestClient_Query_GetMany(t *testing.T) {
	client := testGetClient()

	data, err := client.Query(testGetManyQuery())
	if err != nil {
		log.Fatal()
	}

	// TODO: this will only work on my DB right now
	assert.Equal(
		t,
		map[string][]interface{}{"camera": {
			map[string]interface{}{"name": "Driveway", "stream_url": "rtsp://192.168.137.31:554/Streaming/Channels/101/", "uuid": "3830e9a5-673d-4e7f-ae9b-afa9aeb439ab"},
			map[string]interface{}{"name": "FrontDoor", "stream_url": "rtsp://192.168.137.32:554/Streaming/Channels/101/", "uuid": "cd056389-b0b0-4978-9167-68c93e59f53d"},
			map[string]interface{}{"name": "SideGate", "stream_url": "rtsp://192.168.137.33:554/Streaming/Channels/101/", "uuid": "ba9a4013-9c25-4bd0-931f-6eaf61e7369f"},
		}},
		data,
	)
}

func TestClient_Query_GetOne(t *testing.T) {
	client := testGetClient()

	data, err := client.Query(testGetOneQuery())
	if err != nil {
		log.Fatal()
	}

	// TODO: this will only work on my DB right now
	assert.Equal(
		t,
		map[string][]interface{}{"camera": {map[string]interface{}{"name": "Driveway", "stream_url": "rtsp://192.168.137.31:554/Streaming/Channels/101/", "uuid": "3830e9a5-673d-4e7f-ae9b-afa9aeb439ab"}}},
		data,
	)
}

func TestClient_Extract_GetMany(t *testing.T) {
	client := testGetClient()

	data, err := client.Query(testGetManyQuery())
	if err != nil {
		log.Fatal()
	}

	result := make([]model.Camera, 0)

	err = client.Extract(data, "", &result)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: this will only work on my DB right now
	assert.Equal(
		t,
		[]model.Camera{
			{UUID: result[0].UUID, Name: "Driveway", StreamURL: "rtsp://192.168.137.31:554/Streaming/Channels/101/"},
			{UUID: result[1].UUID, Name: "FrontDoor", StreamURL: "rtsp://192.168.137.32:554/Streaming/Channels/101/"},
			{UUID: result[2].UUID, Name: "SideGate", StreamURL: "rtsp://192.168.137.33:554/Streaming/Channels/101/"}},
		result,
	)
}

func TestClient_QueryAndExtract_GetMany(t *testing.T) {
	client := testGetClient()

	result := make([]model.Camera, 0)

	err := client.QueryAndExtract(
		testGetManyQuery(),
		"",
		&result,
	)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: this will only work on my DB right now
	assert.Equal(
		t,
		[]model.Camera{
			{UUID: result[0].UUID, Name: "Driveway", StreamURL: "rtsp://192.168.137.31:554/Streaming/Channels/101/"},
			{UUID: result[1].UUID, Name: "FrontDoor", StreamURL: "rtsp://192.168.137.32:554/Streaming/Channels/101/"},
			{UUID: result[2].UUID, Name: "SideGate", StreamURL: "rtsp://192.168.137.33:554/Streaming/Channels/101/"}},
		result,
	)
}

func TestClient_QueryAndExtract_InsertOne(t *testing.T) {
	client := testGetClient()

	result := make([]model.Camera, 0)

	err := client.QueryAndExtract(
		testInsertOneQuery(),
		"",
		&result,
	)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(
		t,
		[]model.Camera{
			{UUID: uuid.UUID{0x64, 0xdb, 0xac, 0x5a, 0x29, 0xc7, 0x42, 0x44, 0xb2, 0x97, 0xc, 0x54, 0xa, 0xf3, 0x29, 0xf9}, Name: "TestCamera", StreamURL: "rtsp://192.168.137.34:554/Streaming/Channels/101/"},
		},
		result,
	)

	// error because now there's a duplicate
	err = client.QueryAndExtract(
		testInsertOneQuery(),
		"",
		&result,
	)
	assert.NotNil(t, err)
	assert.Equal(
		t,
		[]model.Camera{
			{UUID: uuid.UUID{0x64, 0xdb, 0xac, 0x5a, 0x29, 0xc7, 0x42, 0x44, 0xb2, 0x97, 0xc, 0x54, 0xa, 0xf3, 0x29, 0xf9}, Name: "TestCamera", StreamURL: "rtsp://192.168.137.34:554/Streaming/Channels/101/"},
		},
		result,
	)

	// okay now delete it
	err = client.QueryAndExtract(
		testDeleteQuery(),
		"",
		&result,
	)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(
		t,
		[]model.Camera{
			{UUID: uuid.UUID{0x64, 0xdb, 0xac, 0x5a, 0x29, 0xc7, 0x42, 0x44, 0xb2, 0x97, 0xc, 0x54, 0xa, 0xf3, 0x29, 0xf9}, Name: "TestCamera", StreamURL: "rtsp://192.168.137.34:554/Streaming/Channels/101/"},
		},
		result,
	)

	// this should work again now
	err = client.QueryAndExtract(
		testInsertOneQuery(),
		"",
		&result,
	)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(
		t,
		[]model.Camera{
			{UUID: uuid.UUID{0x64, 0xdb, 0xac, 0x5a, 0x29, 0xc7, 0x42, 0x44, 0xb2, 0x97, 0xc, 0x54, 0xa, 0xf3, 0x29, 0xf9}, Name: "TestCamera", StreamURL: "rtsp://192.168.137.34:554/Streaming/Channels/101/"},
		},
		result,
	)

	// okay delete it again
	err = client.QueryAndExtract(
		testDeleteQuery(),
		"",
		&result,
	)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(
		t,
		[]model.Camera{
			{UUID: uuid.UUID{0x64, 0xdb, 0xac, 0x5a, 0x29, 0xc7, 0x42, 0x44, 0xb2, 0x97, 0xc, 0x54, 0xa, 0xf3, 0x29, 0xf9}, Name: "TestCamera", StreamURL: "rtsp://192.168.137.34:554/Streaming/Channels/101/"},
		},
		result,
	)
}

func TestClient_QueryAndExtract_MultipleGetMany(t *testing.T) {
	client := testGetClient()

	var err error
	var query string

	for i := 0; i < 4; i++ {
		_, err = client.Query(testInsertOneImage())
		if err != nil {
			log.Fatal(err)
		}

		_, err = client.Query(testInsertOneVideo())
		if err != nil {
			log.Fatal(err)
		}
	}

	images := make([]model.Image, 0)
	videos := make([]model.Video, 0)

	err = client.QueryAndExtractMultiple(
		testGetMultipleManyQuery(),
		[]string{"image", "video"},
		&images,
		&videos,
	)
	if err != nil {
		log.Fatal(err)
	}

	assert.Len(t, images, 4)
	assert.Len(t, videos, 4)

	for _, image := range images {
		query, err = DeleteQuery("image", image)
		if err != nil {
			log.Fatal(err)
		}

		_, err = client.Query(query)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, video := range videos {
		query, err = DeleteQuery("video", video)
		if err != nil {
			log.Fatal(err)
		}

		_, err = client.Query(query)
		if err != nil {
			log.Fatal(err)
		}
	}
}
