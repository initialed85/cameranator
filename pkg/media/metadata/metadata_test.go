package metadata

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetVideoDuration(t *testing.T) {
	duration, err := GetVideoDuration("../../../test_data/segments/Segment_2020-12-25T08:45:04_Driveway.mp4")
	if err != nil {
		require.NoError(t, err)
	}
	assert.Equal(
		t,
		time.Second*30,
		duration,
	)
}

func TestGetSize(t *testing.T) {
	size, err := GetFileSize("../../../test_data/segments/Segment_2020-12-25T08:45:04_Driveway.mp4")
	if err != nil {
		require.NoError(t, err)
	}

	assert.Equal(
		t,
		23.850268,
		size,
	)
}
