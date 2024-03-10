package segment_processor

import (
	"log"
	"strings"
	"time"

	"github.com/initialed85/cameranator/pkg/media/converter"
	"github.com/initialed85/cameranator/pkg/persistence/application"
	"github.com/initialed85/cameranator/pkg/persistence/helpers"
	"github.com/initialed85/cameranator/pkg/segments/event_receiver"
	"github.com/initialed85/cameranator/pkg/segments/segment_generator"
	"github.com/initialed85/cameranator/pkg/utils"
)

type WorkAndError struct {
	Work converter.Work
	Err  error
}

type SegmentProcessor struct {
	correlator     *utils.Correlator
	eventReceiver  *event_receiver.EventReceiver
	imageConverter *converter.Converter
	application    *application.Application
}

func NewSegmentProcessor(
	port int64,
	url string,
	timeout time.Duration,
) (*SegmentProcessor, error) {
	var err error

	m := SegmentProcessor{
		correlator: utils.NewCorrelator(),
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
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (s *SegmentProcessor) eventReceiverHandler(event segment_generator.Event) {
	correlation := s.correlator.NewCorrelation(s.reconcileEvent)

	imageWork := converter.Work{
		SourcePath:      event.ImagePath,
		DestinationPath: strings.ReplaceAll(event.ImagePath, ".jpg", "__lowres.jpg"),
		Width:           640,
		Height:          360,
	}

	imageItem := correlation.NewItem("image")

	eventItem := correlation.NewItem("event")
	eventItem.SetValue(event)
	eventItem.Complete()

	s.imageConverter.Submit(
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

func (s *SegmentProcessor) reconcileEvent(correlation *utils.Correlation) {
	log.Printf("reconciling %#+v...", correlation.GetCorrelationID().String())

	eventItem, err := correlation.GetItem("event")
	if err != nil {
		log.Printf("warning: %#+v marked as complete but failed to get event because %v", correlation, err)
		return
	}

	originalEvent := eventItem.GetValue().(segment_generator.Event)

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
		s.application,
		originalEvent.CameraName,
		originalEvent.VideoStartTimestamp,
		originalEvent.VideoEndTimestamp,
		originalEvent.VideoPath,
		imageWork.Work.DestinationPath,
	)
	if err != nil {
		log.Printf("warning: could not handle event because %v", err)
		return
	}

	log.Printf("added %#+v", event)
}

func (s *SegmentProcessor) Start() error {
	s.imageConverter.Start()
	return s.eventReceiver.Open()
}

func (s *SegmentProcessor) Stop() {
	s.eventReceiver.Close()
	s.imageConverter.Stop()
}
