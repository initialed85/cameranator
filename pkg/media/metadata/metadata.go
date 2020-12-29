package metadata

import (
	"fmt"
	"time"

	"github.com/alfg/mp4"
)

func GetVideoDuration(path string) (time.Duration, error) {
	video, err := mp4.Open(path)
	if err != nil {
		return time.Duration(0), err
	}

	if video.Moov == nil || video.Moov.Mvhd == nil {
		return time.Duration(0), fmt.Errorf("can't get duration- video.Moov or video.Moov.Mvhd is nil")
	}

	return time.Millisecond * time.Duration(video.Moov.Mvhd.Duration), nil
}
