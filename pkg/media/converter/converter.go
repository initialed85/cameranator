package converter

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/initialed85/cameranator/pkg/process"
	"github.com/initialed85/cameranator/pkg/utils"
)

var disableNvidia = false
var enableConversion = false

func DisableNvidia() {
	disableNvidia = true
	log.Printf("warning: Nvidia support disabled at user request")
}

func EnableConversion() {
	enableConversion = true
	log.Printf("warning: conversion support disabled at user request")
}

func init() {
	if os.Getenv("DISABLE_NVIDIA") == "1" {
		DisableNvidia()
	}

	if os.Getenv("ENABLE_CONVERSION") == "1" {
		EnableConversion()
	}
}

func ConvertVideo(sourcePath, destinationPath string, width, height int) (string, string, error) {
	if !enableConversion {
		_ = os.WriteFile(destinationPath, []byte{}, 0644)
		return "warning: conversion support disabled at user request\n", "", nil
	}

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

	return process.RunCommand(
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

	return process.RunCommand(
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
	executor *utils.Executor
	workFn   func(Work) (string, string, error)
}

func NewConverter(numWorkers int, queueSize int, workFn func(Work) (string, string, error)) *Converter {
	c := Converter{
		executor: utils.NewExecutor(numWorkers, queueSize),
		workFn:   workFn,
	}

	return &c
}

func (c *Converter) Submit(work Work, completeFn func(Work, error)) {
	c.executor.Submit(
		func() (interface{}, error) {
			stdout, stderr, err := c.workFn(work)

			if err != nil {
				err = fmt.Errorf("err=%v, stdout=%#+v, stderr=%#+v", err, stdout, stderr)
			}

			return struct{}{}, err
		},
		func(result interface{}, err error) {
			completeFn(work, err)
		},
	)
}

func (c *Converter) Start() {
	c.executor.Start()
}

func (c *Converter) Stop() {
	c.executor.Stop()
}

func NewVideoConverter(numWorkers int, queueSize int) *Converter {
	return NewConverter(
		numWorkers,
		queueSize,
		func(work Work) (string, string, error) {
			return ConvertVideo(
				work.SourcePath,
				work.DestinationPath,
				work.Width,
				work.Height,
			)
		},
	)
}

func NewImageConverter(numWorkers int, queueSize int) *Converter {
	return NewConverter(
		numWorkers,
		queueSize,
		func(work Work) (string, string, error) {
			return ConvertImage(
				work.SourcePath,
				work.DestinationPath,
				work.Width,
				work.Height,
			)
		},
	)
}
