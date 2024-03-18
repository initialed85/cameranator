package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/initialed85/cameranator/pkg/objects/object_tracker"
	"github.com/initialed85/cameranator/pkg/utils"
	"gocv.io/x/gocv"
)

var query = `
query LiveEvents {
	event(where: {status: {_eq: "needs tracking"}, start_timestamp: {_gte: "0000-01-01T00:00:00+00:00"}, original_video: {file_path: {_ilike: "%__segment_file_name__"}}}, order_by: {start_timestamp: desc}, limit: 1) {
	  id
	  original_video {
		id
		file_path
		camera_id
		start_timestamp
		end_timestamp
	  }
	  detections {
		id
		timestamp
		centroid
		bounding_box
		class_id
		class_name
		score
	  }
	}
  }
`

func main() {
	segmentsPath := strings.TrimSpace(os.Getenv("SEGMENTS_PATH"))
	if segmentsPath == "" {
		log.Fatal("SEGMENTS_PATH env empty or unset")
	}

	if len(os.Args) < 2 {
		log.Fatal("first argument must be segment file name")
	}
	segmentFileName := strings.TrimSpace(os.Args[1])
	if segmentFileName == "" {
		log.Fatal("first argument must be segment file name")
	}
	parts := strings.Split(strings.Trim(segmentFileName, "/"), "/")
	segmentFileName = parts[len(parts)-1]

	url := strings.TrimSpace(os.Getenv("URL"))
	if url == "" {
		url = "http://localhost:8080/v1/graphql"
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go utils.WaitForCtrlC(cancel)

	mats := make(chan gocv.Mat) // note: unbuffered

	log.Printf("starting tracker...")
	objectTracker, err := object_tracker.New(segmentsPath, mats)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = objectTracker.Close()
	}()

	log.Printf("starting graphql client...")
	graphqlSubscriptionClient := graphql.NewSubscriptionClient(url)

	handler := func(message []byte, err error) error {
		log.Printf("handling message=%v..., err=%v", string(message)[0:256], err)

		if err != nil {
			log.Printf("handler was invoked with error: %v", err)
			return nil
		}

		if err != nil {
			log.Printf("attempt to read message caused %#+v; ignoring", err)
			return nil
		}

		payload := struct {
			Event []object_tracker.PartialEvent `json:"event"`
		}{}

		err = json.Unmarshal(message, &payload)
		if err != nil {
			log.Printf("attempt to unmarshal message caused %#+v; ignoring", err)
			return nil
		}

		for _, event := range payload.Event {
			err = objectTracker.HandleEvent(event)
			if err != nil {
				log.Printf("attempt to handle event caused: %#+v; ignoring", err)
				continue
			}
		}

		return nil
	}

	go func() {
		time.Sleep(time.Millisecond * 100)

		err := graphqlSubscriptionClient.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()
	defer func() {
		_ = graphqlSubscriptionClient.Close()
	}()

	log.Printf("starting graphql subscription...")
	_, err = graphqlSubscriptionClient.Exec(strings.ReplaceAll(query, "__segment_file_name__", segmentFileName), nil, handler)
	if err != nil {
		log.Fatal(err)
	}

	// note: this has to be created very early on in the life of the program (due to it
	// needing to run on the main thread)
	object_tracker.CreateAndHandleGoCVWindow(ctx, cancel, mats)
}
