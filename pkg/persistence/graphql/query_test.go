package graphql

import (
	"fmt"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/relvacode/iso8601"
	"github.com/stretchr/testify/assert"

	"github.com/initialed85/cameranator/pkg/common"
	"github.com/initialed85/cameranator/pkg/persistence/model"
)

func TestGetManyQuery(t *testing.T) {
	query, err := GetManyQuery("camera", model.Camera{}, "id", "asc")
	if err != nil {
		log.Fatal(err)
	}

	log.Print(query)
	fmt.Println("")

	assert.Equal(
		t,
		`
{
  camera (order_by: {id: asc}) {
    id
    uuid
    name
    stream_url
  }
}
`,
		query,
	)
}

func TestGetManyQuery_Nested(t *testing.T) {
	event := model.Event{

	}

	query, err := GetManyQuery("event", event, "id", "asc")
	if err != nil {
		log.Fatal(err)
	}

	log.Print(query)
	fmt.Println("")

	assert.Equal(
		t,
		`
{
  event (order_by: {id: asc}) {
    id
    uuid
    start_timestamp
    end_timestamp
    is_processed
    high_quality_video_id
    high_quality_video {
      id
      uuid
      start_timestamp
      end_timestamp
      size
      file_path
      is_high_quality
      source_camera_id
      source_camera {
        id
        uuid
        name
        stream_url
      }
    }
    high_quality_image_id
    high_quality_image {
      id
      uuid
      timestamp
      size
      file_path
      is_high_quality
      source_camera_id
      source_camera {
        id
        uuid
        name
        stream_url
      }
    }
    low_quality_video_id
    low_quality_video {
      id
      uuid
      start_timestamp
      end_timestamp
      size
      file_path
      is_high_quality
      source_camera_id
      source_camera {
        id
        uuid
        name
        stream_url
      }
    }
    low_quality_image_id
    low_quality_image {
      id
      uuid
      timestamp
      size
      file_path
      is_high_quality
      source_camera_id
      source_camera {
        id
        uuid
        name
        stream_url
      }
    }
    source_camera_id
    source_camera {
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

func TestGetOneQuery(t *testing.T) {
	query, err := GetOneQuery("camera", model.Camera{}, "uuid", "3830e9a5-673d-4e7f-ae9b-afa9aeb439ab")
	if err != nil {
		log.Fatal(err)
	}

	log.Print(query)
	fmt.Println("")

	assert.Equal(
		t,
		`
{
  camera(where: {uuid: {_eq: "3830e9a5-673d-4e7f-ae9b-afa9aeb439ab"}}, limit: 1, distinct_on: uuid) {
    id
    uuid
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
		UUID:      uuid.UUID{0x64, 0xdb, 0xac, 0x5a, 0x29, 0xc7, 0x42, 0x44, 0xb2, 0x97, 0xc, 0x54, 0xa, 0xf3, 0x29, 0xf9},
		Name:      "model.Camera",
		StreamURL: "rtsp://192.168.137.34:554/Streaming/Channels/101/",
	}

	query, err := InsertQuery("camera", camera)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(query)
	fmt.Println("")

	assert.Equal(
		t,
		`
mutation {
  insert_camera_one(object: {
    name: "model.Camera",
    stream_url: "rtsp://192.168.137.34:554/Streaming/Channels/101/",
    uuid: "64dbac5a-29c7-4244-b297-0c540af329f9"
  }) {
    uuid
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

	image := model.Image{
		UUID:          uuid.UUID{0x64, 0xdb, 0xac, 0x5a, 0x29, 0xc7, 0x42, 0x44, 0xb2, 0x97, 0xc, 0x54, 0xa, 0xf3, 0x29, 0xf9},
		Timestamp:     iso8601.Time{Time: timestamp},
		Size:          65536,
		FilePath:      "/path/to/file",
		IsHighQuality: true,
	}

	query, err := InsertQuery("image", image)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(query)
	fmt.Println("")

	assert.Equal(
		t,
		`
mutation {
  insert_image_one(object: {
    file_path: "/path/to/file",
    is_high_quality: true,
    size: 65536,
    timestamp: "2020-12-26T12:23:54+0930",
    uuid: "64dbac5a-29c7-4244-b297-0c540af329f9"
  }) {
    uuid
    timestamp
    size
    file_path
    is_high_quality
  }
}
`,
		query,
	)
}

func TestInsertQuery_Nested(t *testing.T) {
	camera := model.Camera{
		UUID:      uuid.UUID{0x64, 0xdb, 0xac, 0x5a, 0x29, 0xc7, 0x42, 0x44, 0xb2, 0x97, 0xc, 0x54, 0xa, 0xf3, 0x29, 0xf9},
		Name:      "TestCamera",
		StreamURL: "rtsp://192.168.137.34:554/Streaming/Channels/101/",
	}

	image := model.Image{
		UUID:          uuid.UUID{0x42, 0xed, 0xb, 0xce, 0xd8, 0x94, 0x49, 0xfe, 0xbe, 0xef, 0x7, 0xf5, 0xce, 0x7f, 0xcc, 0xec},
		Timestamp:     utils.GetISO8601Time("2020-03-27T08:30:00+08:00"),
		Size:          65536,
		FilePath:      "/some/path",
		IsHighQuality: true,
		SourceCamera:  camera,
	}

	query, err := InsertQuery("image", image)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(query)
	fmt.Println("")

	assert.Equal(
		t,
		`
mutation {
  insert_image_one(object: {
    file_path: "/some/path",
    is_high_quality: true,
    size: 65536,
    source_camera: {
      data: {
        name: "TestCamera",
        stream_url: "rtsp://192.168.137.34:554/Streaming/Channels/101/",
        uuid: "64dbac5a-29c7-4244-b297-0c540af329f9"
      },
      on_conflict: {
        constraint: camera_uuid_key
        update_columns: [name, stream_url, uuid]
      }
    },
    timestamp: "2020-03-27T08:30:00+0800",
    uuid: "42ed0bce-d894-49fe-beef-07f5ce7fccec"
  }) {
    uuid
    timestamp
    size
    file_path
    is_high_quality
    source_camera {
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

func TestDeleteQuery(t *testing.T) {
	camera := model.Camera{
		UUID:
		uuid.UUID{0x64, 0xdb, 0xac, 0x5a, 0x29, 0xc7, 0x42, 0x44, 0xb2, 0x97, 0xc, 0x54, 0xa, 0xf3, 0x29, 0xf9}, Name: "model.Camera", StreamURL: "rtsp://192.168.137.34:554/Streaming/Channels/101/",
	}

	query, err := DeleteQuery("camera", camera)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(query)
	fmt.Println("")

	assert.Equal(
		t,
		`
mutation {
  delete_camera(where: {name: {_eq: "model.Camera"}, stream_url: {_eq: "rtsp://192.168.137.34:554/Streaming/Channels/101/"}, uuid: {_eq: "64dbac5a-29c7-4244-b297-0c540af329f9"}}) {
    returning {
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
