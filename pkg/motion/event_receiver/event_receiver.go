package event_receiver

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/initialed85/glue/pkg/network"
	"github.com/relvacode/iso8601"
)

type EventID struct {
	CameraID    int64
	CameraName  string
	EventNumber int64
}

type Event struct {
	CameraID            int64
	CameraName          string
	EventNumber         int64
	VideoPath           string
	ImagePath           string
	StartTimestamp      iso8601.Time
	VideoStartTimestamp iso8601.Time
	ImageTimestamp      iso8601.Time
	VideoEndTimestamp   iso8601.Time
	EndTimestamp        iso8601.Time
	Complete            bool
}

type EventReceiver struct {
	mu             sync.Mutex
	eventByEventID map[EventID]Event
	receiver       *network.Receiver
	handler        func(Event)
}

func NewEventReceiver(port int64, handler func(Event)) (*EventReceiver, error) {
	interfaceName, err := network.GetDefaultInterfaceName()
	if err != nil {
		return nil, err
	}

	r := EventReceiver{
		receiver: network.NewReceiver(
			fmt.Sprintf("0.0.0.0:%v", port),
			interfaceName,
		),
		eventByEventID: make(map[EventID]Event),
		handler:        handler,
	}

	return &r, nil
}

func (r *EventReceiver) callback(addr *net.UDPAddr, data []byte) {
	log.Printf("EventReceiver.callback; received: addr=%#+v, data=%#+v)", addr.String(), string(data))

	r.mu.Lock()
	defer r.mu.Unlock()

	parts := strings.Split(
		strings.Trim(
			strings.TrimSpace(string(data)),
			"[]",
		),
		", ",
	)

	receivedTimestamp, err := iso8601.ParseString(
		strings.ReplaceAll(
			parts[0],
			",",
			".",
			),
		)
	if err != nil {
		log.Printf("warning: %v", err)
		return
	}

	cameraID, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		log.Printf("warning: %v", err)
		return
	}

	cameraName := parts[4]

	eventNumber, err := strconv.ParseInt(parts[5], 10, 64)
	if err != nil {
		log.Printf("warning: %v", err)
		return
	}

	eventID := EventID{
		CameraID:    cameraID,
		CameraName:  cameraName,
		EventNumber: eventNumber,
	}

	event, ok := r.eventByEventID[eventID]
	if !ok {
		event = Event{
			CameraID:    cameraID,
			CameraName:  cameraName,
			EventNumber: eventNumber,
		}
	}

	switch parts[1] {
	case "on_event_start":
		event.StartTimestamp = iso8601.Time{Time: receivedTimestamp}
	case "on_movie_start":
		event.VideoStartTimestamp = iso8601.Time{Time: receivedTimestamp}
		event.VideoPath = parts[6]
	case "on_picture_save":
		event.ImageTimestamp = iso8601.Time{Time: receivedTimestamp}
		event.ImagePath = parts[6]
	case "on_movie_end":
		event.VideoEndTimestamp = iso8601.Time{Time: receivedTimestamp}
	case "on_event_end":
		event.EndTimestamp = iso8601.Time{Time: receivedTimestamp}
	}

	r.eventByEventID[eventID] = event

	zeroTimestamp := time.Time{}

	goodStartTimestamp := event.StartTimestamp.After(zeroTimestamp)
	goodVideoStartTimestamp := event.VideoStartTimestamp.After(zeroTimestamp)
	goodImageTimestamp := event.ImageTimestamp.After(zeroTimestamp)
	goodVideoEndTimestamp := event.VideoEndTimestamp.After(zeroTimestamp)
	goodEndTimestamp := event.EndTimestamp.After(zeroTimestamp)

	goodVideoPath := event.VideoPath != ""
	goodImagePath := event.ImagePath != ""

	event.Complete = goodStartTimestamp &&
		goodVideoStartTimestamp &&
		goodImageTimestamp &&
		goodVideoEndTimestamp &&
		goodEndTimestamp &&
		goodVideoPath &&
		goodImagePath

	if event.Complete {
		log.Printf("EventReceiver.callback; complete, invoking handler: event=%#+v", event)
		r.handler(*&event)
		delete(r.eventByEventID, eventID)
	}
}

func (r *EventReceiver) Open() error {
	err := r.receiver.RegisterCallback(r.callback)
	if err != nil {
		log.Fatal(err)
	}

	return r.receiver.Open()
}

func (r *EventReceiver) Close() {
	_ = r.receiver.UnregisterCallback(r.callback)

	r.receiver.Close()
}
