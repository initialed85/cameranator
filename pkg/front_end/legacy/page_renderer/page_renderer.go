package page_renderer

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/initialed85/glue/pkg/worker"

	"github.com/initialed85/cameranator/pkg/front_end/legacy/page_render_utils"
	"github.com/initialed85/cameranator/pkg/persistence/application"
	"github.com/initialed85/cameranator/pkg/persistence/legacy"
)

const (
	fileNamePrefix = "events"
	fileNameSuffix = "html"
)

func getIndexPath(path string) string {
	return filepath.Join(
		path,
		fmt.Sprintf(
			"%v.%v",
			fileNamePrefix,
			fileNameSuffix,
		),
	)
}

func getPagePath(path string, timestamp time.Time) string {
	return filepath.Join(
		path,
		fmt.Sprintf(
			"%v_%v.%v",
			fileNamePrefix,
			timestamp.Format("2006_01_02"),
			fileNameSuffix,
		),
	)
}

func getTitle(isSegment bool) string {
	if isSegment {
		return "Segments"
	}

	return "Events"
}

func truncatePath(path string) string {
	_, file := filepath.Split(path)

	return file
}

func truncatePaths(events []legacy.Event) []legacy.Event {
	newEvents := make([]legacy.Event, 0)

	for _, event := range events {
		newEvents = append(
			newEvents,
			legacy.Event{
				EventID:          event.EventID,
				Timestamp:        event.Timestamp,
				CameraName:       event.CameraName,
				HighResImagePath: truncatePath(event.HighResImagePath),
				LowResImagePath:  truncatePath(event.LowResImagePath),
				HighResVideoPath: truncatePath(event.HighResVideoPath),
				LowResVideoPath:  truncatePath(event.LowResVideoPath),
			},
		)
	}

	return newEvents
}

func cleanFolder(path string) error {
	walkFn := func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		_, file := filepath.Split(path)

		if !strings.HasSuffix(file, fileNameSuffix) {
			return nil
		}

		if !(strings.HasPrefix(file, fileNamePrefix)) {
			return nil
		}

		err = os.Remove(path)
		if err != nil {
			return err
		}

		return nil
	}

	err := filepath.Walk(path, walkFn)
	if err != nil {
		return err
	}

	return nil
}

func writeFile(path, data string) error {
	return ioutil.WriteFile(path, []byte(data), 0644)
}

type PageRenderer struct {
	application     *application.Application
	isSegment       bool
	path            string
	scheduledWorker *worker.ScheduledWorker
	lastEvents      []legacy.Event
}

func NewPageRenderer(
	url string,
	timeout time.Duration,
	isSegment bool,
	path string,
) (*PageRenderer, error) {
	var err error

	p := PageRenderer{
		isSegment:  isSegment,
		path:       path,
		lastEvents: make([]legacy.Event, 0),
	}

	p.application, err = application.NewApplication(url, timeout)
	if err != nil {
		return nil, err
	}

	p.scheduledWorker = worker.NewScheduledWorker(
		func() {},
		p.work,
		func() {},
		time.Second*10,
	)

	return &p, nil
}

func (p *PageRenderer) work() {
	allEvents, err := legacy.GetEvents(p.application, p.isSegment)
	if err != nil {
		log.Printf("warning: failed to get legacy events because %v", err)
		return
	}

	log.Printf("got %v events for processing", len(allEvents))

	if reflect.DeepEqual(allEvents, p.lastEvents) {
		log.Printf("no new events, deferring update (%v and %v)", len(allEvents), len(p.lastEvents))
		return
	}

	eventsByDate := legacy.GetEventsDescendingByDateDescending(allEvents)

	log.Printf("got %v dates to process events for", len(eventsByDate))

	htmlByPath := make(map[string]string)

	now := time.Now()

	for eventsDate, events := range eventsByDate {
		eventsHTML, err := page_render_utils.RenderPage(
			getTitle(p.isSegment),
			truncatePaths(events),
			eventsDate,
			now,
		)
		if err != nil {
			log.Printf("failed to call RenderPage for %v because: %v", eventsDate, err)

			continue
		}

		path := getPagePath(p.path, eventsDate)

		htmlByPath[path] = eventsHTML
	}

	eventsSummaryHTML, err := page_render_utils.RenderSummary(
		getTitle(p.isSegment),
		eventsByDate,
		now,
	)
	if err != nil {
		log.Printf("failed to call RenderSummary because: %v", err)

		return
	}

	path := getIndexPath(p.path)

	htmlByPath[path] = eventsSummaryHTML

	err = cleanFolder(p.path)
	if err != nil {
		log.Printf("warning: failed to call cleanFolder or %#+v because %v", p.path, err)
		return
	}

	for path, html := range htmlByPath {
		err = writeFile(path, html)
		if err != nil {
			log.Printf("warning: failed to write %v bytes to %#+v because %v", len(html), path, err)
			continue
		}
	}

	log.Printf("wrote out %v files", len(htmlByPath))

	p.lastEvents = allEvents
}

func (p *PageRenderer) Start() {
	p.scheduledWorker.Start()
}

func (p *PageRenderer) Stop() {
	p.scheduledWorker.Stop()
}
