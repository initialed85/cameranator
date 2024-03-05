package object_task_scheduler

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/initialed85/cameranator/pkg/persistence/application"
	"github.com/initialed85/cameranator/pkg/persistence/model"
	"github.com/initialed85/glue/pkg/worker"
	"github.com/wagslane/go-rabbitmq"
)

// const tickDuration = time.Second * 600

// const query = `
// query LiveEvents {
// 	event(where: {needs_object_processing: {_eq: true}, is_segment: {_eq: true}, start_timestamp: {_gte: "__timestamp__"}}, order_by: {start_timestamp: desc}) {
//     id
//     high_quality_video {
//       file_path
//       source_camera_id
//       start_timestamp
//       end_timestamp
//     }
//   }
// }
// `

const subscription = `
subscription LiveEvents {
	event(where: {needs_object_processing: {_eq: true}, is_segment: {_eq: true}, start_timestamp: {_gte: "__timestamp__"}}, order_by: {start_timestamp: desc}) {
		id
		high_quality_video {
		file_path
		source_camera_id
		start_timestamp
		end_timestamp
		}
  	}
}
`

const mutation = `
mutation UpdateEvent {
	update_event(where: {id: {_eq: __id__}}, _set: {needs_object_processing: false}) {
		returning {
			id
		}
	}
}
`

const amqpIdentifier = "object_tasks"

type ObjectTaskScheduler struct {
	scheduledWorker           *worker.BlockedWorker
	amqpConn                  *rabbitmq.Conn
	amqpPublisher             *rabbitmq.Publisher
	graphqlSubscriptionClient *graphql.SubscriptionClient
	application               *application.Application
	mu                        *sync.Mutex
	url                       string
	amqp                      string
	skipPublish               bool
}

func NewObjectTaskScheduler(
	url string,
	timeout time.Duration,
	amqp string,
	skipPublish bool,
) (*ObjectTaskScheduler, error) {
	o := ObjectTaskScheduler{
		mu:          new(sync.Mutex),
		url:         url,
		amqp:        amqp,
		skipPublish: skipPublish,
	}

	var err error

	o.application, err = application.NewApplication(url, timeout)
	if err != nil {
		return nil, err
	}

	o.scheduledWorker = worker.NewBlockedWorker(
		o.onStart,
		func() {
			time.Sleep(time.Second * 1)
		},
		o.onStop,
	)

	return &o, nil
}

func (o *ObjectTaskScheduler) handler(message []byte, err error) error {
	// message should be the bytes for something like this:
	/*
		{
		  "uuid": "363530a1-0666-4587-a44a-a45bcc644779",
		  "start_timestamp": "2023-01-29T04:35:00+00:00",
		  "end_timestamp": "2023-01-29T04:40:00+00:00",
		  "high_quality_video": {
			"file_path": "/srv/target_dir/segments/Segment_2023-01-29T12:35:00_FrontDoor.mp4"
		  }
		}
	*/

	log.Printf("handling message=%v", string(message))

	if err != nil {
		log.Printf("attempt to read message caused %#+v; ignoring", err)
		return nil
	}

	payload := struct {
		Event []PartialEvent `json:"event"`
	}{}

	err = json.Unmarshal(message, &payload)
	if err != nil {
		log.Printf("attempt to unmarshal message caused %#+v; ignoring", err)
		return nil
	}

	eventModelAndClient, err := o.application.GetModelAndClient("event")
	if err != nil {
		log.Printf("attempt to get model and client caused %#+v; ignoring", err)
		return nil
	}

	for _, event := range payload.Event {
		client := eventModelAndClient.Client()

		updatedEvents := make([]*model.Event, 0)
		err = eventModelAndClient.GetOne(&updatedEvents, "id", event.ID)
		if err != nil {
			log.Printf("attempt to get latest event caused %#+v; ignoring", err)
			continue
		}
		if len(updatedEvents) == 0 {
			log.Printf("attempt to get latest event resulted in an empty set; ignoring")
			continue
		}

		updatedEvent := updatedEvents[0]
		if !updatedEvent.NeedsObjectProcessing {
			log.Printf("event %v has already been processed; skipping", updatedEvent.ID)
			continue
		}

		mutation := strings.ReplaceAll(mutation, "__id__", fmt.Sprintf("%v", event.ID))
		log.Printf("mutation: %v", mutation)

		result, err := client.Mutate(mutation)
		if err != nil {
			log.Printf("attempt to run mutation caused %#+v; ignoring", err)
			continue
		}
		log.Printf("result: %#+v", result)

		eventJSON, err := json.Marshal(event)
		if err != nil {
			log.Printf("attempt to marshal event caused %#+v; ignoring", err)
			continue
		}

		if o.skipPublish {
			log.Printf("skipped publishing event=%v", string(eventJSON))
			continue
		}

		o.mu.Lock()
		err = o.amqpPublisher.Publish(eventJSON, []string{amqpIdentifier})
		o.mu.Unlock()
		if err != nil {
			log.Printf("attempt to publish event caused %#+v; ignoring", err)
			continue
		}

		log.Printf("published event=%v", string(eventJSON))
	}

	return nil
}

func (o *ObjectTaskScheduler) onStart() {
	var err error

	log.Printf("connecting to %v", o.amqp)
	o.amqpConn, err = rabbitmq.NewConn(
		o.amqp,
		rabbitmq.WithConnectionOptionsLogging,
		rabbitmq.WithConnectionOptionsReconnectInterval(time.Second),
	)
	if err != nil {
		// TODO
		log.Panicf("attempt to invoke rabbitmq.NewConn caused %#+v; cannot recover", err)
		return
	}

	log.Printf("creating publisher for %v", amqpIdentifier)
	o.amqpPublisher, err = rabbitmq.NewPublisher(
		o.amqpConn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(amqpIdentifier),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
		rabbitmq.WithPublisherOptionsExchangeAutoDelete,
	)
	if err != nil {
		// TODO
		log.Panicf("attempt to invoke rabbitmq.NewPublisher caused %#+v; cannot recover", err)
		return
	}

	log.Printf("connecting to %v", o.url)
	o.graphqlSubscriptionClient = graphql.NewSubscriptionClient(o.url)

	log.Printf("building subscription...")
	// timestamp := time.Now().UTC().Format(time.RFC3339)
	timestamp := time.Time{}.Format(time.RFC3339) // unix epoch (so, forever)
	subscription := strings.ReplaceAll(subscription, "__timestamp__", timestamp)
	_, err = o.graphqlSubscriptionClient.Exec(subscription, nil, o.handler)
	if err != nil {
		// TODO
		log.Panicf("attempt to invoke graphqlClient.Exec (for subscription) caused %#+v; cannot recover", err)
		return
	}

	log.Printf("%v", subscription)

	log.Printf("running graphql client...")
	err = o.graphqlSubscriptionClient.Run()
	if err != nil {
		// TODO
		log.Panicf("attempt to invoke graphqlClient.Run caused %#+v; cannot recover", err)
		return
	}
}

func (o *ObjectTaskScheduler) onStop() {
	_ = o.graphqlSubscriptionClient.Close()
	o.graphqlSubscriptionClient = nil

	_ = o.amqpPublisher.Close
	o.amqpPublisher = nil

	_ = o.amqpConn.Close()
	o.amqpConn = nil
}

func (o *ObjectTaskScheduler) Start() {
	o.scheduledWorker.Start()
}

func (o *ObjectTaskScheduler) Stop() {
	o.scheduledWorker.Stop()
}
