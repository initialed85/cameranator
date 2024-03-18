package object_tracker

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/initialed85/cameranator/pkg/persistence/application"
	"github.com/initialed85/glue/pkg/worker"
	"gocv.io/x/gocv"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// TODO: fix hacked file_path match
const subscription = `
subscription LiveEvents {
	event(where: {status: {_eq: "needs tracking"}, start_timestamp: {_gte: "__timestamp__"}, original_video: {file_path: {_eq: "/srv/target_dir/segments/Segment_2024-03-11T15:56:01_Driveway.mp4"}}}, order_by: {start_timestamp: desc}, limit: 1) {
	  id
	  original_video {
		file_path
		camera_id
		start_timestamp
		end_timestamp
	  }
	  detections {
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

const mutation = `
mutation UpdateEvent {
	update_event(where: {id: {_in: [__ids__]}}, _set: {status: "tracking underway"}) {
		returning {
			id
		}
	}
}
`

type ObjectTracker struct {
	scheduledWorker           *worker.BlockedWorker
	graphqlSubscriptionClient *graphql.SubscriptionClient
	application               *application.Application
	mu                        *sync.Mutex
	url                       string
	mats                      chan gocv.Mat
}

func NewObjectTracker(
	url string,
	timeout time.Duration,
	mats chan gocv.Mat,
) (*ObjectTracker, error) {
	o := ObjectTracker{
		mu:   new(sync.Mutex),
		url:  url,
		mats: mats,
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

func (o *ObjectTracker) handleEvent(event *PartialEvent) error {
	for i, detection := range event.PartialDetection {
		rawXYAsStr := strings.ReplaceAll(detection.RawCentroid, "(", "")
		rawXYAsStr = strings.ReplaceAll(rawXYAsStr, ")", "")

		rawXYAsSlice := strings.Split(strings.Trim(rawXYAsStr, " ()"), ",")

		x, err := strconv.ParseFloat(rawXYAsSlice[0], 64)
		if err != nil {
			log.Printf("warning: failed to parse %#+v[0] as float", rawXYAsSlice)
			continue
		}

		y, err := strconv.ParseFloat(rawXYAsSlice[1], 64)
		if err != nil {
			log.Printf("warning: failed to parse %#+v[1] as float", rawXYAsSlice)
			continue
		}

		detection.Centroid = Point{X: x, Y: y}

		detection.BoundingBox = make([]Point, 0)
		for _, rawXYAsStr := range strings.Split(detection.RawBoundingBox, "),(") {
			rawXYAsStr = strings.ReplaceAll(rawXYAsStr, "(", "")
			rawXYAsStr = strings.ReplaceAll(rawXYAsStr, ")", "")
			rawXYAsSlice := strings.Split(strings.Trim(rawXYAsStr, " ()"), ",")

			x, err := strconv.ParseFloat(rawXYAsSlice[0], 64)
			if err != nil {
				log.Printf("warning: failed to parse %#+v[0] as float", rawXYAsSlice)
				continue
			}

			y, err := strconv.ParseFloat(rawXYAsSlice[1], 64)
			if err != nil {
				log.Printf("warning: failed to parse %#+v[1] as float", rawXYAsSlice)
				continue
			}

			detection.BoundingBox = append(detection.BoundingBox, Point{X: x, Y: y})
		}
		event.PartialDetection[i] = detection
	}

	filePath := "test_data/segments/Segment_2024-03-11T15_56_01_Driveway.mp4"

	rawData, err := ffmpeg.Probe(filePath)
	if err != nil {
		return fmt.Errorf("failed to probe %#+v", filePath)
	}

	data := []byte(rawData)

	var ffProbeOutput FFProbeOutput

	err = json.Unmarshal(data, &ffProbeOutput)
	if err != nil {
		return fmt.Errorf("failed to parse ffprobe output: %v", err)
	}

	width := ffProbeOutput.Streams[0].Width
	height := ffProbeOutput.Streams[0].Height

	frameRate, err := strconv.ParseInt(strings.Split(ffProbeOutput.Streams[0].RawFrameRate, "/")[0], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse left portion of %#+v as int", ffProbeOutput.Streams[0].RawFrameRate)
	}
	_ = frameRate

	durationSeconds, err := strconv.ParseFloat(ffProbeOutput.Streams[0].RawDuration, 64)
	if err != nil {
		return fmt.Errorf("failed to parse %#+v as int", ffProbeOutput.Streams[0].RawDuration)
	}

	duration := time.Second * time.Duration(durationSeconds)
	_ = duration

	reader, writer := io.Pipe()

	go func() {
		err := ffmpeg.Input(filePath).
			Output("pipe:",
				ffmpeg.KwArgs{
					"format":  "rawvideo",
					"pix_fmt": "rgb24",
				}).
			WithOutput(writer).
			// ErrorToStdOut().
			Run()
		if err != nil {
			log.Panicf("failed to read from %#+v using ffmpeg: %v", filePath, err)
		}
	}()

	frameSize := width * height * 3
	buf := make([]byte, frameSize)

	red := color.RGBA{
		R: uint8(255),
		G: uint8(0),
		B: uint8(0),
	}

	green := color.RGBA{
		R: uint8(0),
		G: uint8(255),
		B: uint8(0),
	}

	white := color.RGBA{
		R: uint8(255),
		G: uint8(255),
		B: uint8(255),
	}

	frame := 1

	for {
		n, err := io.ReadFull(reader, buf)
		if n != frameSize || (err != nil && err != io.EOF) {
			return fmt.Errorf("failed to read %#+v after %v bytes: %v", filePath, n, err)
		}

		if n == 0 || err == io.EOF {
			break
		}

		originalImage, err := gocv.NewMatFromBytes(height, width, gocv.MatTypeCV8UC3, buf)
		if err != nil {
			return fmt.Errorf("failed to decode %#+v bytes: %v", n, err)
		}

		if originalImage.Empty() {
			break
		}

		frameDurationSeconds := float64(frame) / float64(frameRate)
		frameDuration := time.Nanosecond * time.Duration(frameDurationSeconds*1000000000)
		frameTimestamp := event.OriginalVideo.StartTimestamp.Add(frameDuration)

		frame++

		overlay := gocv.NewMat()

		gocv.CvtColor(originalImage, &overlay, gocv.ColorBGRToRGB)
		originalImage.Close()

		for _, detection := range event.PartialDetection {
			if detection.Timestamp.After(frameTimestamp.Add(-time.Millisecond*3000)) &&
				detection.Timestamp.Before(frameTimestamp) {

				gocv.Circle(
					&overlay,
					image.Point{int(detection.Centroid.X), int(detection.Centroid.Y)},
					10,
					white,
					1,
				)
			}

			if detection.Timestamp.After(frameTimestamp.Add(-time.Millisecond*250)) &&
				detection.Timestamp.Before(frameTimestamp.Add(+time.Millisecond*250)) {
				topLeft := detection.BoundingBox[0]
				bottomRight := detection.BoundingBox[2]

				r := image.Rectangle{
					Min: image.Point{
						X: int(topLeft.X),
						Y: int(topLeft.Y),
					},
					Max: image.Point{
						X: int(bottomRight.X),
						Y: int(bottomRight.Y),
					},
				}

				gocv.Rectangle(&overlay, r, green, 1)

				text := fmt.Sprintf("%v (%.2f%%)", detection.ClassName, detection.Score*100.0)
				textSize := gocv.GetTextSize(text, gocv.FontHersheyPlain, 1.5, 1)
				pt := image.Pt(int(detection.Centroid.X)-(textSize.X/2), int(detection.Centroid.Y))
				gocv.PutText(&overlay, text, pt, gocv.FontHersheyPlain, 2.0, red, 2)
			}
		}

		o.mats <- overlay
	}

	log.Printf("done")

	return nil
}

func (o *ObjectTracker) handler(message []byte, err error) error {
	if err != nil {
		return err
	}

	log.Printf("handling message=%v, err=%v", string(message), err)

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

	wg := new(sync.WaitGroup)

	for _, event := range payload.Event {
		event := event
		wg.Add(1)
		go func() {
			defer wg.Done()

			err = o.handleEvent(&event)
			if err != nil {
				log.Printf("warning: failed to handle event: %v", err)
				return
			}
		}()
	}

	wg.Wait()

	// eventIDs := make([]string, 0)
	// for _, event := range payload.Event {
	// 	eventIDs = append(eventIDs, fmt.Sprintf("%v", event.ID))
	// }

	// eventModelAndClient, err := o.application.GetModelAndClient("event")
	// if err != nil {
	// 	log.Printf("attempt to get model and client caused %#+v; ignoring", err)
	// 	return nil
	// }

	// client := eventModelAndClient.Client()

	// mutation := strings.ReplaceAll(mutation, "__ids__", strings.Join(eventIDs, ", "))
	// log.Printf("mutation: %v", mutation)

	// result, err := client.Mutate(mutation)
	// if err != nil {
	// 	log.Printf("attempt to run mutation caused %#+v; ignoring", err)
	// 	return nil
	// }
	// log.Printf("result: %#+v", result)

	return nil
}

func (o *ObjectTracker) onStart() {
	var err error

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

func (o *ObjectTracker) onStop() {
	_ = o.graphqlSubscriptionClient.Close()
	o.graphqlSubscriptionClient = nil
}

func (o *ObjectTracker) Start() {
	o.scheduledWorker.Start()
}

func (o *ObjectTracker) Stop() {
	o.scheduledWorker.Stop()
}
