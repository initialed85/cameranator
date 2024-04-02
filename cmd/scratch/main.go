package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/initialed85/cameranator/pkg/objects/object_tracker"
	"github.com/initialed85/cameranator/pkg/utils"
	"gocv.io/x/gocv"
)

// getQuery 2024-03-30T06:17:01+00:00
func getQuery(cameraName string, startTimestamp time.Time) string {
	query := strings.TrimSpace(
		fmt.Sprintf(`
query LiveEvents {
	event(where: {source_camera: {name: {_eq: "%v"}}, start_timestamp: {_gte: "%v"}, status: {_eq: "needs tracking"}}, limit: 1, order_by: {start_timestamp: asc}) {
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
}`,
			cameraName,
			startTimestamp.Format(time.RFC3339),
		),
	)

	return query
}

func main() {
	segmentsPath := strings.TrimSpace(os.Getenv("SEGMENTS_PATH"))
	if segmentsPath == "" {
		log.Fatal("SEGMENTS_PATH env empty or unset")
	}

	if len(os.Args) < 3 {
		log.Fatal("first argument must be camera name, second argument must by start timestamp")
	}

	cameraName := strings.TrimSpace(os.Args[1])
	if cameraName == "" {
		log.Fatal("first argument must be camera name")
	}
	parts := strings.Split(strings.Trim(cameraName, "/"), "/")
	cameraName = parts[len(parts)-1]

	rawStartTimestamp := strings.TrimSpace(os.Args[2])
	if rawStartTimestamp == "" {
		log.Fatal("second argument must be start timestamp")
	}
	startTimestamp, err := time.Parse(time.RFC3339, rawStartTimestamp)
	if err != nil {
		log.Fatalf("second argument must be start timestamp: %v", err)
	}

	url := strings.TrimSpace(os.Getenv("URL"))
	if url == "" {
		url = "https://cameranator.initialed85.cc/api/v1/graphql"
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

	log.Printf("creating graphql client...")
	graphqlClient := graphql.NewClient(url, http.DefaultClient)

	events := make(chan object_tracker.PartialEvent, 1)

	go func() {
		for {
			log.Printf("waiting for an event to handle...")

			event := <-events

			log.Printf("handling event %v", event.ID)

			err = objectTracker.HandleEvent(event)
			if err != nil {
				log.Printf("attempt to handle event caused: %#+v; ignoring", err)
				continue
			}

			log.Printf("handled event %v", event.ID)
		}
	}()

	go func() {
		timestamp := startTimestamp

		for {
			log.Printf("querying for cameraName: %v, startTimestamp: %v...", cameraName, timestamp.Format(time.RFC3339))

			before := time.Now()

			b, err := graphqlClient.ExecRaw(context.Background(), getQuery(cameraName, timestamp), nil)
			if err != nil {
				log.Fatal(err)
			}

			after := time.Now()

			log.Printf("queried for cameraName: %v, startTimestamp: %v in %v", cameraName, timestamp.Format(time.RFC3339), after.Sub(before))

			payload := struct {
				Event []object_tracker.PartialEvent `json:"event"`
			}{}

			err = json.Unmarshal(b, &payload)
			if err != nil {
				log.Printf("attempt to unmarshal message caused %#+v; ignoring", err)
				continue
			}

			for _, event := range payload.Event {
				log.Printf("waiting to queue event %v...", event.ID)

				events <- event

				log.Printf("queued event %v", event.ID)

				timestamp = event.OriginalVideo.EndTimestamp.Time
			}
		}
	}()

	// note: this has to be created very early on in the life of the program (due to it
	// needing to run on the main thread)
	log.Printf("opening window...")
	object_tracker.CreateAndHandleGoCVWindow(ctx, cancel, mats)
}
