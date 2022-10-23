package segment_generator

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/relvacode/iso8601"

	"github.com/initialed85/cameranator/pkg/filesystem"
	"github.com/initialed85/cameranator/pkg/media/metadata"
	"github.com/initialed85/cameranator/pkg/media/segment_recorder"
	"github.com/initialed85/cameranator/pkg/media/thumbnail_creator"
	"github.com/initialed85/cameranator/pkg/process"
	"github.com/initialed85/cameranator/pkg/utils"
)

type Event struct {
	CameraName          string
	VideoPath           string
	ImagePath           string
	VideoStartTimestamp iso8601.Time
	VideoEndTimestamp   iso8601.Time
	ImageTimestamp      iso8601.Time
}

type Feed struct {
	NetCamURL       string
	DestinationPath string
	CameraName      string
	Duration        int
}

type SegmentGenerator struct {
	feed                 Feed
	completeFn           func(event Event)
	mu                   sync.Mutex
	backgroundProcess    *process.BackgroundProcess
	watcher              *filesystem.Watcher
	lastCreatedPath      string
	lastCreatedTimestamp time.Time
}

func NewSegmentGenerator(
	feed Feed,
	completeFn func(Event),
) *SegmentGenerator {
	s := SegmentGenerator{
		feed:                 feed,
		completeFn:           completeFn,
		lastCreatedTimestamp: time.Now(),
	}

	return &s
}

func (s *SegmentGenerator) onFileCreate(file filesystem.File) {
	s.mu.Lock()
	lastCreatedPath := s.lastCreatedPath
	s.mu.Unlock()

	if file.Name == lastCreatedPath {
		return
	}

	if lastCreatedPath != "" {
		log.Printf("onFileCreate; %#+v closed, %#+v created", lastCreatedPath, file.Name)

		imagePath := fmt.Sprintf("%v.jpg", strings.Split(lastCreatedPath, ".mp4")[0])

		err := thumbnail_creator.GetThumbnail(
			lastCreatedPath,
			imagePath,
		)
		if err != nil {
			log.Printf("warning: attempt to get thumbnail for %#+v raisd %#+v", lastCreatedPath, err)
		}

		// TODO: nasty hack for local timezone
		rawVideoStartTimestamp := fmt.Sprintf(
			"%v+08:00",
			strings.Split(
				strings.Split(
					lastCreatedPath,
					"/Segment_",
				)[1],
				fmt.Sprintf("_%v", s.feed.CameraName),
			)[0])

		videoStartTimestamp := utils.GetISO8601Time(rawVideoStartTimestamp)

		duration, err := metadata.GetVideoDuration(lastCreatedPath)
		if err != nil {
			log.Printf("warning: attempt to get duration for %#+v raised %#+v", lastCreatedPath, err)
			return
		}

		videoEndTimestamp := iso8601.Time{Time: videoStartTimestamp.Add(duration)}

		imageTimestamp := videoStartTimestamp

		event := Event{
			CameraName:          s.feed.CameraName,
			VideoPath:           lastCreatedPath,
			ImagePath:           imagePath,
			VideoStartTimestamp: videoStartTimestamp,
			VideoEndTimestamp:   videoEndTimestamp,
			ImageTimestamp:      imageTimestamp,
		}

		log.Printf("onFileCreate; complete event=%#+v", event)

		s.completeFn(event)

		s.mu.Lock()
		s.lastCreatedTimestamp = time.Now()
		s.mu.Unlock()
	}

	s.mu.Lock()
	s.lastCreatedPath = file.Name
	s.mu.Unlock()
}

func (s *SegmentGenerator) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	backgroundProcess, err := segment_recorder.RecordSegments(
		s.feed.NetCamURL,
		s.feed.DestinationPath,
		s.feed.CameraName,
		s.feed.Duration,
	)
	if err != nil {
		return err
	}

	s.backgroundProcess = backgroundProcess

	// TODO: depends on not changing the filename pattern in RecordSegments
	matcher, err := regexp.Compile(".*/Segment_.*_" + s.feed.CameraName + "\\.mp4")
	if err != nil {
		return err
	}

	s.watcher = filesystem.NewWatcher(
		s.feed.DestinationPath,
		matcher,
		s.onFileCreate,
		func(file filesystem.File) {},
	)

	s.watcher.Start()

	return nil
}

func (s *SegmentGenerator) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.backgroundProcess.Stop()
	s.watcher.Stop()
}

func (s *SegmentGenerator) IsLive() bool {
	s.mu.Lock()
	lastCreatedTimestamp := s.lastCreatedTimestamp
	s.mu.Unlock()

	expiry := lastCreatedTimestamp.Add(time.Second * time.Duration(float64(s.feed.Duration)*1.5))

	return time.Now().Before(expiry)
}
