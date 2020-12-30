package segment_recorder

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/initialed85/cameranator/pkg/process"
)

var disableNvidia = false

func DisableNvidia() {
	disableNvidia = true
	log.Printf("warning: Nvidia support disabled at user request")
}

func init() {
	if os.Getenv("DISABLE_NVIDIA") == "1" {
		DisableNvidia()
	}
}

func RecordSegments(netCamURL, destinationPath, cameraName string, duration int) (*process.BackgroundProcess, error) {
	log.Printf("RecordSegments; recording %v second segments from %v to %v for %v", duration, netCamURL, destinationPath, cameraName)

	arguments := make([]string, 0)

	if !disableNvidia {
		arguments = append(
			arguments,
			"-hwaccel",
			"cuda",
			"-c:v",
			"h264_cuvid",
		)
	} else {
		arguments = append(
			arguments,
			"-c:v",
			"h264",
		)
	}

	arguments = append(
		arguments,
		"-rtsp_transport",
		"tcp",
		"-i",
		netCamURL,
		"-c",
		"copy",
		"-map",
		"0",
		"-f",
		"segment",
		"-segment_time",
		fmt.Sprintf("%v", duration),
		"-segment_format",
		"mp4",
		"-segment_atclocktime",
		"1",
		"-strftime",
		"1",
		"-x264-params",
		"keyint=100:scenecut=0",
		"-g",
		"100",
		"-muxdelay",
		"0",
		"-muxpreload",
		"0",
		"-reset_timestamps",
		"1",
	)

	if !disableNvidia {
		arguments = append(
			arguments,
			"-c:v",
			"h264_nvenc",
		)
	} else {
		arguments = append(
			arguments,
			"-c:v",
			"libx264",
		)
	}

	arguments = append(
		arguments,
		filepath.Join(destinationPath, "Segment_%Y-%m-%dT%H:%M:%S_"+cameraName+".mp4"),
	)

	return process.RunBackgroundProcess(
		"ffmpeg",
		arguments...,
	)
}
