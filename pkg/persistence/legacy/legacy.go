package legacy

import (
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/initialed85/cameranator/pkg/persistence/application"
	"github.com/initialed85/cameranator/pkg/persistence/model"
)

type Event struct {
	EventID          uuid.UUID `json:"event_id"`
	Timestamp        time.Time `json:"timestamp"`
	CameraName       string    `json:"camera_name"`
	HighResImagePath string    `json:"high_res_image_path"`
	LowResImagePath  string    `json:"low_res_image_path"`
	HighResVideoPath string    `json:"high_res_video_path"`
	LowResVideoPath  string    `json:"low_res_video_path"`
}

func GetEvents(application *application.Application, isSegment bool) ([]Event, error) {
	eventModelAndClient, err := application.GetModelAndClient("event")
	if err != nil {
		return []Event{}, err
	}

	events := make([]model.Event, 0)
	err = eventModelAndClient.GetAll(&events)
	if err != nil {
		return []Event{}, err
	}

	legacyEvents := make([]Event, 0)
	for _, event := range events {
		// TODO: better served by a where clause
		if event.IsSegment != isSegment {
			continue
		}

		legacyEvent := Event{
			EventID:          event.UUID,
			Timestamp:        event.StartTimestamp.Time,
			CameraName:       event.SourceCamera.Name,
			HighResImagePath: event.HighQualityImage.FilePath,
			LowResImagePath:  event.LowQualityImage.FilePath,
			HighResVideoPath: event.HighQualityVideo.FilePath,
			LowResVideoPath:  event.LowQualityVideo.FilePath,
		}

		legacyEvents = append(legacyEvents, legacyEvent)
	}

	return legacyEvents, nil
}

func GetEventsDescending(events []Event) []Event {
	sort.SliceStable(events, func(i, j int) bool {
		return events[i].Timestamp.Unix() > events[j].Timestamp.Unix()
	})

	return events
}

func GetEventsDescendingByDateDescending(events []Event) map[time.Time][]Event {
	allEvents := GetEventsDescending(events)

	eventsByDate := make(map[time.Time][]Event)

	for _, event := range allEvents {
		date, _ := time.Parse("2006-01-02", event.Timestamp.Format("2006-01-02"))

		_, ok := eventsByDate[date]
		if !ok {
			eventsByDate[date] = make([]Event, 0)
		}

		eventsByDate[date] = append(eventsByDate[date], event)
	}

	for date := range eventsByDate {
		sort.SliceStable(eventsByDate[date], func(i, j int) bool {
			return eventsByDate[date][i].Timestamp.Unix() > eventsByDate[date][j].Timestamp.Unix()
		})
	}

	keys := make([]time.Time, 0)
	for key := range eventsByDate {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i].Unix() > keys[j].Unix()
	})

	sortedEventsByDate := make(map[time.Time][]Event)
	for _, key := range keys {
		sortedEventsByDate[key] = eventsByDate[key]
	}

	return sortedEventsByDate
}
