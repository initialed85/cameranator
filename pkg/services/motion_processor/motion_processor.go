package motion_processor

import (
	"log"
	"strings"
	"time"

	"github.com/initialed85/cameranator/pkg/media/converter"
	"github.com/initialed85/cameranator/pkg/motion/event_receiver"
	"github.com/initialed85/cameranator/pkg/persistence/application"
	"github.com/initialed85/cameranator/pkg/persistence/helpers"
	"github.com/initialed85/cameranator/pkg/utils"
)

type WorkAndError struct {
	Work converter.Work
	Err  error
}

type MotionProcessor struct {
	correlator     *utils.Correlator
	eventReceiver  *event_receiver.EventReceiver
	videoConverter *converter.Converter
	imageConverter *converter.Converter
	application    *application.Application
}

func NewMotionProcessor(
	port int64,
	url string,
	timeout time.Duration,
) (*MotionProcessor, error) {
	var err error

	m := MotionProcessor{
		correlator: utils.NewCorrelator(),
		videoConverter: converter.NewVideoConverter(
			2,
			1024,
		),
		imageConverter: converter.NewImageConverter(
			2,
			1024,
		),
	}

	m.eventReceiver, err = event_receiver.NewEventReceiver(port, m.eventReceiverHandler)
	if err != nil {
		return nil, err
	}

	m.application, err = application.NewApplication(url, timeout)

	return &m, nil
}

func (m *MotionProcessor) eventReceiverHandler(event event_receiver.Event) {
	correlation := m.correlator.NewCorrelation(m.reconcileEvent)

	videoWork := converter.Work{
		SourcePath:      event.VideoPath,
		DestinationPath: strings.ReplaceAll(event.VideoPath, ".mp4", "__lowres.mp4"),
		Width:           640,
		Height:          360,
	}

	imageWork := converter.Work{
		SourcePath:      event.ImagePath,
		DestinationPath: strings.ReplaceAll(event.ImagePath, ".jpg", "__lowres.jpg"),
		Width:           640,
		Height:          360,
	}

	eventItem := correlation.NewItem("event")
	videoItem := correlation.NewItem("video")
	imageItem := correlation.NewItem("image")

	eventItem.SetValue(event)
	eventItem.Complete()

	m.videoConverter.Submit(
		videoWork,
		func(work converter.Work, err error) {
			videoItem.SetValue(WorkAndError{
				Work: work,
				Err:  err,
			})
			videoItem.Complete()
		},
	)

	m.imageConverter.Submit(
		imageWork,
		func(work converter.Work, err error) {
			imageItem.SetValue(WorkAndError{
				Work: work,
				Err:  err,
			})
			imageItem.Complete()
		},
	)
}

func (m *MotionProcessor) reconcileEvent(correlation *utils.Correlation) {
	log.Printf("reconciling %#+v...", correlation.GetCorrelationID().String())

	eventItem, err := correlation.GetItem("event")
	if err != nil {
		log.Printf("warning: %#+v marked as complete but failed to get event because %v", correlation, err)
		return
	}

	originalEvent := eventItem.GetValue().(event_receiver.Event)

	videoWorkItem, err := correlation.GetItem("video")
	if err != nil {
		log.Printf("warning: %#+v marked as complete but failed to get video because %v", correlation, err)
		return
	}

	videoWork := videoWorkItem.GetValue().(WorkAndError)
	if videoWork.Err != nil {
		log.Printf("warning: %#+v marked as complete but failed to get video because %v", correlation, videoWork.Err)
		return
	}

	imageWorkItem, err := correlation.GetItem("image")
	if err != nil {
		log.Printf("warning: %#+v marked as complete but failed to get image because %v", correlation, err)
		return
	}

	imageWork := imageWorkItem.GetValue().(WorkAndError)
	if imageWork.Err != nil {
		log.Printf("warning: %#+v marked as complete but failed to get image because %v", correlation, imageWork.Err)
		return
	}

	event, err := helpers.AddEvent(
		m.application,
		originalEvent.CameraName,
		originalEvent.StartTimestamp,
		originalEvent.EndTimestamp,
		false,
		videoWork.Work.SourcePath,
		imageWork.Work.SourcePath,
		videoWork.Work.DestinationPath,
		imageWork.Work.DestinationPath,
	)
	if err != nil {
		log.Printf("warning: could not handle event because %v", err)
		return
	}

	log.Printf("added %#+v", event)
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
