package graphql

import (
	"testing"

	"github.com/relvacode/iso8601"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/initialed85/cameranator/pkg/persistence/model"
	"github.com/initialed85/cameranator/pkg/utils"
)

func TestGetManyQuery(t *testing.T) {
	query, err := GetManyQuery("camera", model.Camera{}, "", nil, "id", "asc")
	require.NoError(t, err)

	assert.Equal(
		t,
		`
{
  camera (order_by: {id: asc}) {
    id
    name
    stream_url
  }
}
`,
		query,
	)
}

func TestGetManyQuery_Nested(t *testing.T) {
	event := model.Event{}

	query, err := GetManyQuery("event", event, "", nil, "id", "asc")
	require.NoError(t, err)

	assert.Equal(
		t,
		`
{
  event (order_by: {id: asc}) {
    id
    start_timestamp
    end_timestamp
    original_video_id
    original_video {
      id
      start_timestamp
      end_timestamp
      size
      file_path
      camera_id
      camera {
        id
        name
        stream_url
      }
    }
    thumbnail_image_id
    thumbnail_image {
      id
      timestamp
      size
      file_path
      camera_id
      camera {
        id
        name
        stream_url
      }
    }
    source_camera_id
    source_camera {
      id
      name
      stream_url
    }
    status
  }
}
`,
		query,
	)
}

func TestGetOneQuery(t *testing.T) {
	query, err := GetOneQuery("camera", model.Camera{}, "name", "Driveway")
	require.NoError(t, err)

	assert.Equal(
		t,
		`
{
  camera(where: {name: {_eq: "Driveway"}}, limit: 1, distinct_on: name) {
    id
    name
    stream_url
  }
}
`,
		query,
	)
}

func TestInsertQuery(t *testing.T) {
	camera := model.Camera{
		Name:      "Driveway1",
		StreamURL: "rtsp://192.168.137.31:554/Streaming/Channels/101/",
	}

	query, err := InsertQuery("camera", camera)
	require.NoError(t, err)

	assert.Equal(
		t,
		`
mutation {
  insert_camera_one(object: {
    name: "Driveway1",
    stream_url: "rtsp://192.168.137.31:554/Streaming/Channels/101/"
  }) {
    id
    name
    stream_url
  }
}
`,
		query,
	)
}

func TestInsertQuery_WithTimestamp(t *testing.T) {
	timestamp, _ := iso8601.ParseString("2020-12-26T12:23:54+0930")

	camera := model.Camera{
		Name:      "Driveway2",
		StreamURL: "rtsp://192.168.137.31:554/Streaming/Channels/101/",
	}

	image := model.Image{
		Timestamp: iso8601.Time{Time: timestamp},
		Size:      65536,
		FilePath:  "/path/to/file",
		Camera:    camera,
	}

	query, err := InsertQuery("image", image)
	require.NoError(t, err)

	assert.Equal(
		t,
		`
mutation {
  insert_image_one(object: {
    camera: {
      data: {
        name: "Driveway2",
        stream_url: "rtsp://192.168.137.31:554/Streaming/Channels/101/"
      },
      on_conflict: {
        constraint: camera_pkey
        update_columns: [name, stream_url]
      }
    },
    file_path: "/path/to/file",
    size: 65536,
    timestamp: "2020-12-26T12:23:54+0930"
  }) {
    id
    timestamp
    size
    file_path
    camera_id
    camera {
      id
      name
      stream_url
    }
  }
}
`,
		query,
	)
}

func TestInsertQuery_Nested(t *testing.T) {
	camera := model.Camera{
		Name:      "Driveway3",
		StreamURL: "rtsp://192.168.137.31:554/Streaming/Channels/101/",
	}

	image := model.Image{
		Timestamp: utils.GetISO8601Time("2020-03-27T08:30:00+08:00"),
		Size:      65536,
		FilePath:  "/some/path",
		Camera:    camera,
	}

	query, err := InsertQuery("image", image)
	require.NoError(t, err)

	assert.Equal(
		t,
		`
mutation {
  insert_image_one(object: {
    camera: {
      data: {
        name: "Driveway3",
        stream_url: "rtsp://192.168.137.31:554/Streaming/Channels/101/"
      },
      on_conflict: {
        constraint: camera_pkey
        update_columns: [name, stream_url]
      }
    },
    file_path: "/some/path",
    size: 65536,
    timestamp: "2020-03-27T08:30:00+0800"
  }) {
    id
    timestamp
    size
    file_path
    camera_id
    camera {
      id
      name
      stream_url
    }
  }
}
`,
		query,
	)
}

func TestDeleteQuery(t *testing.T) {
	camera := model.Camera{
		Name:      "Driveway4",
		StreamURL: "rtsp://192.168.137.31:554/Streaming/Channels/101/",
	}

	query, err := DeleteQuery("camera", camera)
	require.NoError(t, err)

	assert.Equal(
		t,
		`
mutation {
  delete_camera(where: {name: {_eq: "Driveway4"}, stream_url: {_eq: "rtsp://192.168.137.31:554/Streaming/Channels/101/"}}) {
    returning {
      id
      name
      stream_url
    }
  }
}
`,
		query,
	)
}
