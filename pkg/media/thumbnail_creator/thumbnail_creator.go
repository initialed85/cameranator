package thumbnail_creator

import (
	"fmt"

	"github.com/initialed85/cameranator/pkg/process"
)

func GetThumbnail(videoPath, imagePath string) error {
	stdout, stderr, err := process.RunCommand(
		"ffmpeg",
		"-i",
		videoPath,
		"-ss",
		"00:00:00.000",
		"-vframes",
		"1",
		imagePath,
	)
	if err != nil {
		return fmt.Errorf("%v; stdout=%#+v, stderr=%#+v", err, stdout, stderr)
	}

	return nil
}
