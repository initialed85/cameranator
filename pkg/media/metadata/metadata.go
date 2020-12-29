package metadata

import (
	"fmt"
	"os"
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

func GetFileSize(path string) (float64, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}

	if fileInfo.IsDir() {
		return 0, fmt.Errorf("%#+v is a folder, not a file", path)
	}

	return float64(fileInfo.Size()) / 1000000, nil
}
