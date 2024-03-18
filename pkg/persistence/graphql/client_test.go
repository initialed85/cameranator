package graphql

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/initialed85/cameranator/pkg/persistence/model"
)

func testGetClient() *Client {
	return NewClient("http://localhost:8082/v1/graphql", time.Second*30)
}

func testInsertOneQuery() string {
	return `
mutation {
  insert_camera_one(object: {
		name: "TestCamera_TestClient",
		stream_url: "rtsp://192.168.137.34:554/Streaming/Channels/101/"
	}) {
    name
    stream_url
  }
}
`
}

func testInsertOneWithOnConflictQuery() string {
	return `
mutation {
  insert_camera_one(object: {
		name: "TestCamera_TestClient",
		stream_url: "rtsp://192.168.137.34:554/Streaming/Channels/101/"
	},
	on_conflict: {
		constraint: camera_name_key,
		update_columns: name
	}) {
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
    file_path: "path/to/file",
    camera: {
		data: {
			name: "TestCamera_TestClient",
			stream_url: "rtsp://192.168.137.34:554/Streaming/Channels/101/"
		},
		on_conflict: {
			constraint: camera_name_key,
			update_columns: name
		}
	}
  }) {
    id
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
    file_path: "path/to/file",
    camera: {
		data: {
			name: "TestCamera_TestClient",
			stream_url: "rtsp://192.168.137.34:554/Streaming/Channels/101/"
			},
			on_conflict: {
				constraint: camera_name_key,
				update_columns: name
			}
		}
  }) {
    id
  }
}
`
}

func testGetManyQuery() string {
	return `
{
  camera (order_by: {id: asc}) {
    name
    stream_url
  }
}
	`
}

func testGetOneQuery() string {
	return `
{
  camera(where: {name: {_eq: "TestCamera_TestClient"}}, limit: 1, distinct_on: name) {
    name
    stream_url
  }
}
`
}

func testDeleteQuery() string {
	return `
mutation {
  delete_camera(where: {name: {_eq: "TestCamera_TestClient"}}) {
    returning {
      id
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
    timestamp
    size
    file_path
    camera {
      name
      stream_url
    }
  }
  video (order_by: {id: asc}) {
    start_timestamp
    end_timestamp
    size
    file_path
    camera {
      name
      stream_url
    }
  }
}
	`
}

func TestClient_QueryAndExtract_InsertOne(t *testing.T) {
	client := testGetClient()

	result := make([]model.Camera, 0)

	err := client.QueryAndExtract(
		testInsertOneQuery(),
		"",
		&result,
	)
	require.NoError(t, err)
	assert.Equal(
		t,
		[]model.Camera{
			{Name: "TestCamera_TestClient", StreamURL: "rtsp://192.168.137.34:554/Streaming/Channels/101/"},
		},
		result,
	)

	// error because now there's a duplicate
	err = client.QueryAndExtract(
		testInsertOneQuery(),
		"",
		&result,
	)
	require.Error(t, err)
	assert.Equal(
		t,
		[]model.Camera{
			{Name: "TestCamera_TestClient", StreamURL: "rtsp://192.168.137.34:554/Streaming/Channels/101/"},
		},
		result,
	)

	// okay now delete it
	err = client.QueryAndExtract(
		testDeleteQuery(),
		"",
		&result,
	)
	require.NoError(t, err)
	assert.Equal(
		t,
		[]model.Camera{
			{Name: "TestCamera_TestClient", StreamURL: "rtsp://192.168.137.34:554/Streaming/Channels/101/"},
		},
		result,
	)

	// this should work again now
	err = client.QueryAndExtract(
		testInsertOneQuery(),
		"",
		&result,
	)
	require.NoError(t, err)
	assert.Equal(
		t,
		[]model.Camera{
			{Name: "TestCamera_TestClient", StreamURL: "rtsp://192.168.137.34:554/Streaming/Channels/101/"},
		},
		result,
	)
}

func TestClient_Query_GetMany(t *testing.T) {
	client := testGetClient()

	data, err := client.Query(testGetManyQuery())
	require.NoError(t, err)

	assert.Condition(t, func() bool {
		for _, rawItem := range data["camera"] {
			item := rawItem.(map[string]any)
			if item["name"] == "TestCamera_TestClient" && item["stream_url"] == "rtsp://192.168.137.34:554/Streaming/Channels/101/" {
				return true
			}
		}
		return false
	})
}

func TestClient_Query_GetOne(t *testing.T) {
	client := testGetClient()

	data, err := client.Query(testGetOneQuery())
	require.NoError(t, err)

	// TODO: this will only work on my DB right now
	assert.Equal(
		t,
		map[string][]interface{}{"camera": {
			map[string]interface{}{"name": "TestCamera_TestClient", "stream_url": "rtsp://192.168.137.34:554/Streaming/Channels/101/"}},
		},
		data,
	)
}

func TestClient_Extract_GetMany(t *testing.T) {
	client := testGetClient()

	data, err := client.Query(testGetManyQuery())
	require.NoError(t, err)

	result := make([]model.Camera, 0)

	err = client.Extract(data, "", &result)
	require.NoError(t, err)

	assert.Condition(t, func() bool {
		for _, camera := range result {
			if camera.Name == "TestCamera_TestClient" && camera.StreamURL == "rtsp://192.168.137.34:554/Streaming/Channels/101/" {
				return true
			}
		}
		return false
	})
}

func TestClient_QueryAndExtract_GetMany(t *testing.T) {
	client := testGetClient()

	result := make([]model.Camera, 0)

	err := client.QueryAndExtract(
		testGetManyQuery(),
		"",
		&result,
	)
	require.NoError(t, err)

	assert.Condition(t, func() bool {
		for _, camera := range result {
			if camera.Name == "TestCamera_TestClient" && camera.StreamURL == "rtsp://192.168.137.34:554/Streaming/Channels/101/" {
				return true
			}
		}
		return false
	})
}

func TestClient_QueryAndExtract_MultipleGetMany(t *testing.T) {
	client := testGetClient()

	var err error
	var query string

	for i := 0; i < 4; i++ {
		_, err = client.Query(testInsertOneImage())
		if err != nil {
			require.NoError(t, err)
		}

		_, err = client.Query(testInsertOneVideo())
		if err != nil {
			require.NoError(t, err)
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
	require.NoError(t, err)

	assert.Len(t, images, 4)
	assert.Len(t, videos, 4)

	for _, image := range images {
		query, err = DeleteQuery("image", image)
		if err != nil {
			require.NoError(t, err)
		}

		_, err = client.Query(query)
		if err != nil {
			require.NoError(t, err)
		}
	}

	for _, video := range videos {
		query, err = DeleteQuery("video", video)
		if err != nil {
			require.NoError(t, err)
		}

		_, err = client.Query(query)
		if err != nil {
			require.NoError(t, err)
		}
	}
}
