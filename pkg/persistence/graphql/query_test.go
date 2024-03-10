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
	query, err := GetManyQuery("camera", model.Camera{}, "id", "asc")
	if err != nil {
		require.NoError(t, err)
	}

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

	query, err := GetManyQuery("event", event, "id", "asc")
	if err != nil {
		require.NoError(t, err)
	}

	assert.Equal(
		t,
		`
{
  event (order_by: {id: asc}) {
    id
    start_timestamp
    end_timestamp
    is_segment
    video_id
    video {
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
    low_quality_video_id
    low_quality_video {
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
    low_quality_image_id
    low_quality_image {
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

func TestGetOneQuery(t *testing.T) {
	query, err := GetOneQuery("camera", model.Camera{}, "name", "Driveway")
	if err != nil {
		require.NoError(t, err)
	}

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
	if err != nil {
		require.NoError(t, err)
	}

	assert.Equal(
		t,
		`
mutation {
  insert_camera_one(object: {
    name: "Driveway1",
    stream_url: "rtsp://192.168.137.31:554/Streaming/Channels/101/",
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
	if err != nil {
		require.NoError(t, err)
	}

	assert.Equal(
		t,
		`
mutation {
  insert_image_one(object: {
    file_path: "/path/to/file",
    size: 65536,
    camera: {
      data: {
        name: "Driveway2",
        stream_url: "rtsp://192.168.137.31:554/Streaming/Channels/101/",
      },
      on_conflict: {
        constraint: camera_name_key
        update_columns: [name, stream_url]
      }
    },
    timestamp: "2020-12-26T12:23:54+0930",
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
	if err != nil {
		require.NoError(t, err)
	}

	assert.Equal(
		t,
		`
mutation {
  insert_image_one(object: {
    file_path: "/some/path",
    size: 65536,
    camera: {
      data: {
        name: "Driveway3",
        stream_url: "rtsp://192.168.137.31:554/Streaming/Channels/101/",
      },
      on_conflict: {
        constraint: camera_id_key
        update_columns: [name, stream_url]
      }
    },
    timestamp: "2020-03-27T08:30:00+0800",
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
	if err != nil {
		require.NoError(t, err)
	}

	assert.Equal(
		t,
		`
mutation {
  delete_camera(where: {name: {_eq: "Driveway4"}, stream_url: {_eq: "rtsp://192.168.137.31:554/Streaming/Channels/101/"}}) {
    returning {
      id
      uuid
      name
      stream_url
    }
  }
}
`,
		query,
	)
}
