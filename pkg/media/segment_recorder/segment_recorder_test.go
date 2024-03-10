package segment_recorder

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/relvacode/iso8601"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecordSegments(t *testing.T) {
	dir, err := os.MkdirTemp("", "cameranator")
	if err != nil {
		require.NoError(t, err)
	}

	process, err := RecordSegments("rtsp://host.docker.internal:8554/Streaming/Channels/101", dir, "Driveway", 5)
	if process == nil {
		require.Fail(t, "process unexpectedly nil")
	}

	if err != nil {
		process.Stop()

		require.NoError(t, err)
	}

	time.Sleep(time.Second * 10)

	process.Stop()

	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		require.NoError(t, err)
	}

	assert.Greater(t, len(fileInfos), 0)

	for _, fileInfo := range fileInfos {
		assert.True(t, strings.HasPrefix(fileInfo.Name(), "Segment_"))
		assert.True(t, strings.HasSuffix(fileInfo.Name(), "_Driveway.mp4"))
		timestamp, err := iso8601.ParseString(
			strings.Split(fileInfo.Name(), "_")[1],
		)
		if err != nil {
			require.NoError(t, err)
		}
		log.Printf("parsed: %#+v", timestamp.String())
	}
}
