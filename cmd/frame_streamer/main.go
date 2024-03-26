package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/initialed85/cameranator/pkg/media/display"
	"github.com/initialed85/cameranator/pkg/media/frame"
	"github.com/initialed85/cameranator/pkg/persistence/model"
	"github.com/initialed85/cameranator/pkg/utils"
	"github.com/relvacode/iso8601"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"github.com/vmihailenco/msgpack/v5"
	"gocv.io/x/gocv"

	"github.com/initialed85/glue/pkg/endpoint"
)

func main() {
	rawCameraID := strings.TrimSpace(os.Getenv("CAMERA_ID"))
	if rawCameraID == "" {
		log.Fatal("CAMERA_ID env var empty or unset")
	}
	cameraID, err := strconv.ParseInt(rawCameraID, 10, 64)
	if err != nil {
		log.Fatalf("CAMERA_ID could not be parsed as int: %v", err)
	}
	if cameraID <= 0 {
		log.Fatal("CAMERA_ID less than or equal to 0")
	}

	cameraName := strings.TrimSpace(os.Getenv("CAMERA_NAME"))
	if cameraName == "" {
		log.Fatal("CAMERA_NAME env var empty or unset")
	}

	streamURL := strings.TrimSpace(os.Getenv("STREAM_URL"))
	if streamURL == "" {
		log.Fatal("STREAM_URL env var empty or unset")
	}

	endpointManager, err := endpoint.NewManagerSimple()
	if err != nil {
		log.Fatal(err)
	}

	endpointManager.Start()
	defer endpointManager.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go utils.WaitForCtrlC(cancel)

	reader, writer := io.Pipe()
	defer func() {
		_ = reader.Close()
		_ = writer.Close()
	}()

	go func() {
		<-ctx.Done()
		_ = reader.Close()
		_ = writer.Close()
	}()

	frames := make(chan frame.Frame, 60*20) // 1 minute at 20 fps
	mats := make(chan gocv.Mat, 60*20)      // 1 minute at 20 fps
	errs := make(chan error, 16)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case frame := <-frames:
				b, err := msgpack.Marshal(frame)
				if err != nil {
					log.Printf("warning: %v", err)
					continue
				}

				err = endpointManager.Publish(
					fmt.Sprintf("raw_frames/%v", cameraID),
					"RawFrame",
					time.Minute,
					b,
				)
				if err != nil {
					log.Printf("warning: %v", err)
				}
			}
		}
	}()

	go func() {
		ffmpeg := ffmpeg_go.Input(
			streamURL,
			ffmpeg_go.KwArgs{
				"rtsp_transport": "tcp",
			},
		).
			Output(
				"pipe:",
				ffmpeg_go.KwArgs{
					"format":  "rawvideo",
					"pix_fmt": "rgb24",
				},
			).
			WithOutput(writer)

		if os.Getenv("DEBUG") == "1" {
			ffmpeg = ffmpeg.ErrorToStdOut()
		}

		err := ffmpeg.Run()
		if err != nil {
			errs <- err
		}
	}()

	go func() {
		select {
		case <-time.After(time.Second * 5):
		case err := <-errs:
			log.Fatalf("failed to open ffmpeg input stream: %v: %v", streamURL, err)
		}

		width := int(1920)  // TODO
		height := int(1080) // TODO
		frameSize := width * height * 3
		buf := make([]byte, frameSize)

		for {
			n, err := io.ReadFull(reader, buf)
			if n != frameSize || (err != nil && err != io.EOF) {
				log.Fatalf("failed to read %#+v after %v bytes: %v", streamURL, n, err)
			}

			now := time.Now()

			if n == 0 || err == io.EOF {
				log.Fatal("stream stopped sending")
			}

			originalImage, err := gocv.NewMatFromBytes(height, width, gocv.MatTypeCV8UC3, buf)
			if err != nil {
				log.Fatalf("failed to decode %#+v bytes: %v", n, err)
			}

			if originalImage.Empty() {
				log.Fatal("original image was empty")
			}

			frame := frame.Frame{
				Camera: model.Camera{
					ID:        cameraID,
					Name:      cameraName,
					StreamURL: streamURL,
				},
				Timestamp: iso8601.Time{Time: now},
				Data:      buf,
			}

			log.Printf("buf: %v", len(buf))

			select {
			case frames <- frame:
			default:
				log.Printf("warning: glue publish delays causing dropped frames")
			}

			overlay := gocv.NewMat()
			gocv.CvtColor(originalImage, &overlay, gocv.ColorBGRToRGBA)
			originalImage.Close()

			if mats != nil {
				select {
				case mats <- overlay:
				default:
					log.Printf("warning: window delays causing dropped frames")
				}
			}
		}
	}()

	if os.Getenv("DEBUG") == "1" {
		log.Printf("debug enabled, showing window...")

		err := display.CreateAndHandleGoCVWindow(ctx, cancel, mats)
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	<-ctx.Done()
}
