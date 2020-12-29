package metadata

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetVideoDuration(t *testing.T) {
	duration, err := GetVideoDuration("../../../test_data/segments/Segment_2020-12-25_08-45-04_Driveway.mp4")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(
		t,
		time.Minute*5,
		duration,
	)
}
