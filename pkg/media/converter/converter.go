package converter

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/initialed85/glue/pkg/worker"

	"github.com/initialed85/cameranator/pkg/common"
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

func ConvertVideo(sourcePath, destinationPath string, width, height int) (string, string, error) {
	var err error

	sourcePath, err = filepath.Abs(sourcePath)
	if err != nil {
		return "", "", err
	}

	destinationPath, err = filepath.Abs(destinationPath)
	if err != nil {
		return "", "", err
	}

	sourcePath = strings.TrimSpace(sourcePath)
	destinationPath = strings.TrimSpace(destinationPath)

	log.Printf("ConvertVideo; sourcePath=%#+v, destinationPath=%#+v, width=%#+v, height=%#+v", sourcePath, destinationPath, width, height)

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
		"-y",
		"-i",
		sourcePath,
		"-s",
		fmt.Sprintf("%vx%v", width, height),
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
		destinationPath,
	)

	return common.RunCommand(
		"ffmpeg",
		arguments...,
	)
}

func ConvertImage(sourcePath, destinationPath string, width, height int) (string, string, error) {
	var err error

	sourcePath, err = filepath.Abs(sourcePath)
	if err != nil {
		return "", "", err
	}

	destinationPath, err = filepath.Abs(destinationPath)
	if err != nil {
		return "", "", err
	}

	sourcePath = strings.TrimSpace(sourcePath)
	destinationPath = strings.TrimSpace(destinationPath)

	log.Printf("ConvertImage; sourcePath=%#+v, destinationPath=%#+v, width=%#+v, height=%#+v", sourcePath, destinationPath, width, height)

	arguments := []string{
		"-resize",
		fmt.Sprintf("%vX%v", width, height),
		sourcePath,
		destinationPath,
	}

	return common.RunCommand(
		"convert",
		arguments...,
	)
}

type Work struct {
	SourcePath      string
	DestinationPath string
	Width           int
	Height          int
}

type Converter struct {
	blockedWorker *worker.BlockedWorker
	workQueue     chan Work
	workFn        func(string, string, int, int) (string, string, error)
	completeFn    func(Work, string, string, error)
}

func NewConverter(
	workQueue chan Work,
	workFn func(string, string, int, int) (string, string, error),
	completeFn func(Work, string, string, error),
) *Converter {
	c := Converter{
		workQueue:  workQueue,
		workFn:     workFn,
		completeFn: completeFn,
	}

	c.blockedWorker = worker.NewBlockedWorker(
		func() {},
		func() {
			c.work()
		},
		func() {
			close(c.workQueue)
		},
	)

	return &c
}

func (c *Converter) work() {
	log.Printf("Converter.work; waiting for work...")

	work := <-c.workQueue
	log.Printf("Converter.work; got %#+v, working...", work)

	stdout, stderr, err := c.workFn(
		work.SourcePath,
		work.DestinationPath,
		work.Width,
		work.Height,
	)

	log.Printf("Converter.work; calling completeFn with...\nstdout=%v\n\nstderr=%v\n\nerr=%#+v", stdout, stderr, err)
	c.completeFn(work, stdout, stderr, err)
}

func (c *Converter) Start() {
	c.blockedWorker.Start()
}

func (c *Converter) Stop() {
	c.blockedWorker.Stop()
}

func NewVideoConverter(
	workQueue chan Work,
	completeFn func(Work, string, string, error),
) *Converter {
	return NewConverter(
		workQueue,
		ConvertVideo,
		completeFn,
	)
}

func NewImageConverter(
	workQueue chan Work,
	completeFn func(Work, string, string, error),
) *Converter {
	return NewConverter(
		workQueue,
		ConvertImage,
		completeFn,
	)
}
