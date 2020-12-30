package motion_processor

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/initialed85/cameranator/pkg/media/converter"
	"github.com/initialed85/cameranator/pkg/motion/event_receiver"
	"github.com/initialed85/cameranator/pkg/persistence/application"
	"github.com/initialed85/cameranator/pkg/persistence/helpers"
)

type Conversion struct {
	VideoWork converter.Work
	ImageWork converter.Work
}

type ConvertedEvent struct {
	Event         event_receiver.Event
	VideoComplete bool
	ImageComplete bool
}

type MotionProcessor struct {
	mu                         sync.Mutex
	eventReceiver              *event_receiver.EventReceiver
	convertedEventByConversion map[Conversion]*ConvertedEvent
	videoConverterWorkQueue    chan converter.Work
	videoConverter             *converter.Converter
	imageConverterWorkQueue    chan converter.Work
	imageConverter             *converter.Converter
	application                *application.Application
}

func NewMotionProcessor(
	port int64,
	url string,
	timeout time.Duration,
) (*MotionProcessor, error) {
	var err error

	m := MotionProcessor{
		convertedEventByConversion: make(map[Conversion]*ConvertedEvent),
		videoConverterWorkQueue:    make(chan converter.Work, 1024),
		imageConverterWorkQueue:    make(chan converter.Work, 1024),
	}

	m.eventReceiver, err = event_receiver.NewEventReceiver(port, m.eventReceiverHandler)
	if err != nil {
		return nil, err
	}

	m.videoConverter = converter.NewVideoConverter(
		m.videoConverterWorkQueue,
		m.videoConverterCompleteFn,
	)

	m.imageConverter = converter.NewImageConverter(
		m.imageConverterWorkQueue,
		m.imageConverterCompleteFn,
	)

	m.application, err = application.NewApplication(url, timeout)

	return &m, nil
}

func (m *MotionProcessor) eventReceiverHandler(event event_receiver.Event) {
	conversion := Conversion{
		VideoWork: converter.Work{
			SourcePath:      event.VideoPath,
			DestinationPath: strings.ReplaceAll(event.VideoPath, ".mp4", "__lowres.mp4"),
			Width:           640,
			Height:          360,
		},
		ImageWork: converter.Work{
			SourcePath:      event.ImagePath,
			DestinationPath: strings.ReplaceAll(event.ImagePath, ".jpg", "__lowres.jpg"),
			Width:           640,
			Height:          360,
		},
	}

	m.videoConverterWorkQueue <- conversion.VideoWork

	m.imageConverterWorkQueue <- conversion.ImageWork

	m.mu.Lock()
	defer m.mu.Unlock()

	m.convertedEventByConversion[conversion] = &ConvertedEvent{
		Event:         event,
		VideoComplete: false,
		ImageComplete: false,
	}
}

func (m *MotionProcessor) getConvertedEvent(work converter.Work) (*ConvertedEvent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for conversion, convertedEvent := range m.convertedEventByConversion {
		if !(conversion.VideoWork == work || conversion.ImageWork == work) {
			continue
		}

		return convertedEvent, nil
	}

	return nil, fmt.Errorf("failed to find convertedEvent for %#+v", work)
}

func (m *MotionProcessor) reconcileConvertedEvent(convertedEvent *ConvertedEvent) {
	if !(convertedEvent.VideoComplete && convertedEvent.ImageComplete) {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	found := false
	var conversion Conversion
	var possibleConvertedEvent *ConvertedEvent

	for conversion, possibleConvertedEvent = range m.convertedEventByConversion {
		found = true

		if convertedEvent.Event.VideoPath != possibleConvertedEvent.Event.VideoPath {
			found = false
		}

		if convertedEvent.Event.ImagePath != possibleConvertedEvent.Event.ImagePath {
			found = false
		}

		if !found {
			continue
		}

		break
	}

	if !found {
		log.Printf(
			"warning: could not handle event because %v",
			fmt.Errorf("failed to find conversion for convertedEvent"),
		)
		return
	}

	delete(m.convertedEventByConversion, conversion)

	event, err := helpers.AddEvent(
		m.application,
		convertedEvent.Event.CameraName,
		convertedEvent.Event.StartTimestamp,
		convertedEvent.Event.EndTimestamp,
		conversion.VideoWork.SourcePath,
		conversion.ImageWork.SourcePath,
		conversion.VideoWork.DestinationPath,
		conversion.ImageWork.DestinationPath,
	)
	if err != nil {
		log.Printf("warning: could not handle event because %v", err)
		return
	}

	log.Printf("added %#+v", event)
}

func (m *MotionProcessor) videoConverterCompleteFn(work converter.Work, stdout string, stderr string, err error) {
	convertedEvent, findConvertedEventErr := m.getConvertedEvent(work)
	if findConvertedEventErr != nil {
		log.Printf("warning: could not convert video because %v", findConvertedEventErr)
		return
	}

	if err != nil {
		log.Printf("warning:could not convert video because %v", err)
		return
	}

	convertedEvent.VideoComplete = true

	m.reconcileConvertedEvent(convertedEvent)
}

func (m *MotionProcessor) imageConverterCompleteFn(work converter.Work, stdout string, stderr string, err error) {
	convertedEvent, findConvertedEventErr := m.getConvertedEvent(work)
	if findConvertedEventErr != nil {
		log.Printf("warning: could not convert image because %v", findConvertedEventErr)
		return
	}

	if err != nil {
		log.Printf("warning:could not convert image because %v", err)
		return
	}

	convertedEvent.ImageComplete = true

	m.reconcileConvertedEvent(convertedEvent)
}

func (m *MotionProcessor) Start() error {
	m.imageConverter.Start()
	m.videoConverter.Start()
	return m.eventReceiver.Open()
}

func (m *MotionProcessor) Stop() {
	m.eventReceiver.Close()
	m.videoConverter.Stop()
	m.imageConverter.Stop()
}
