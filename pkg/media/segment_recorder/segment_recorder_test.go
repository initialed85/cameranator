package segment_recorder

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/relvacode/iso8601"
	"github.com/stretchr/testify/assert"
)

func TestRecordSegments(t *testing.T) {
	DisableNvidia()

	dir, err := ioutil.TempDir("", "cameranator")
	if err != nil {
		log.Fatal(err)
	}

	process, err := RecordSegments("rtsp://localhost:8554/Streaming/Channels/101", dir, "Driveway", 5)
	if process == nil {
		log.Fatal("process unexpectedly nil")
	}

	if err != nil {
		process.Stop()

		log.Fatal(err)
	}

	time.Sleep(time.Second * 10)

	process.Stop()

	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	assert.Greater(t, len(fileInfos), 0)

	for _, fileInfo := range fileInfos {
		assert.True(t, strings.HasPrefix(fileInfo.Name(), "Segment_"))
		assert.True(t, strings.HasSuffix(fileInfo.Name(), "_Driveway.mp4"))
		timestamp, err := iso8601.ParseString(
			strings.Split(fileInfo.Name(), "_")[1],
		)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("parsed: %#+v", timestamp.String())
	}
}
