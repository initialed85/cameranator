package event_receiver

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/initialed85/glue/pkg/network"

	"github.com/initialed85/cameranator/pkg/segments/segment_generator"
)

type EventReceiver struct {
	mu       sync.Mutex
	receiver *network.Receiver
	handler  func(segment_generator.Event)
}

func NewEventReceiver(port int64, handler func(event segment_generator.Event)) (*EventReceiver, error) {
	interfaceName, err := network.GetDefaultInterfaceName()
	if err != nil {
		return nil, err
	}

	r := EventReceiver{
		receiver: network.NewReceiver(
			fmt.Sprintf("0.0.0.0:%v", port),
			interfaceName,
		),
		handler: handler,
	}

	return &r, nil
}

func (r *EventReceiver) callback(addr *net.UDPAddr, data []byte) {
	log.Printf("EventReceiver.callback; received: addr=%#+v, data=%#+v)", addr.String(), string(data))

	r.mu.Lock()
	defer r.mu.Unlock()

	event := segment_generator.Event{}
	err := json.Unmarshal(data, &event)
	if err != nil {
		log.Printf("warning: attempt to unmarshal %#+v raised %v", string(data), err)
		return
	}

	log.Printf("EventReceiver.callback; complete, invoking handler: event=%#+v", event)
	r.handler(event)
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
