package object_task_scheduler

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/hasura/go-graphql-client"
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

const amqpIdentifier = "object_tasks"

type ObjectTaskScheduler struct {
	scheduledWorker *worker.BlockedWorker
	amqpConn        *rabbitmq.Conn
	amqpPublisher   *rabbitmq.Publisher
	graphqlClient   *graphql.SubscriptionClient
	url             string
	amqp            string
}

func NewObjectTaskScheduler(
	url string,
	amqp string,
) (*ObjectTaskScheduler, error) {
	o := ObjectTaskScheduler{
		url:  url,
		amqp: amqp,
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

	for _, event := range payload.Event {
		eventJSON, err := json.Marshal(event)
		if err != nil {
			log.Printf("attempt to marshal event caused %#+v; ignoring", err)
			continue
		}

		err = o.amqpPublisher.Publish(eventJSON, []string{amqpIdentifier})
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

	o.graphqlClient = graphql.NewSubscriptionClient(o.url)

	subscription := strings.ReplaceAll(subscription, "__timestamp__", time.Now().UTC().Format(time.RFC3339))
	_, err = o.graphqlClient.Exec(subscription, nil, o.handler)
	if err != nil {
		// TODO
		log.Panicf("attempt to invoke graphqlClient.Exec (for subscription) caused %#+v; cannot recover", err)
		return
	}

	err = o.graphqlClient.Run()
	if err != nil {
		// TODO
		log.Panicf("attempt to invoke graphqlClient.Run caused %#+v; cannot recover", err)
		return
	}
}

func (o *ObjectTaskScheduler) onStop() {
	_ = o.graphqlClient.Close()
	o.graphqlClient = nil

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
