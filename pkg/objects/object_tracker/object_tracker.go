package object_tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"math"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gocv.io/x/gocv"
	"golang.org/x/exp/maps"
)

const (
	objectDetectionStrideFrames = 4
	boundingBoxHoldDuration     = time.Millisecond * 1_00
	centroidHoldDuration        = time.Millisecond * 10_000
	objectLookbackDuration      = time.Millisecond * 10_000
	areaFactorLimit             = 0.25
	aspectRatioFactorLimit      = 0.25
	centroidDistanceFactorLimit = 0.25
)

type ObjectTracker struct {
	segmentsPath string
	mats         chan gocv.Mat
	ctx          context.Context
	cancel       context.CancelFunc
}

func New(
	segmentsPath string,
	mats chan gocv.Mat,
) (*ObjectTracker, error) {
	segmentsPath = strings.TrimRight(segmentsPath, "/")

	o := ObjectTracker{
		segmentsPath: segmentsPath,
		mats:         mats,
	}
	o.ctx, o.cancel = context.WithCancel(context.Background())

	return &o, nil
}

func (o *ObjectTracker) Close() error {
	o.cancel()
	return nil
}

func (o *ObjectTracker) enrichEvent(event *PartialEvent) error {
	log.Printf("event.ID: %v; preparing original video %#+v...", event.ID, event.OriginalVideo.FilePath)

	filePathParts := strings.Split(strings.TrimRight(event.OriginalVideo.FilePath, "/"), "/")
	event.OriginalVideo.AdjustedFilePath = path.Join(o.segmentsPath, filePathParts[len(filePathParts)-1])
	log.Printf("event.ID: %v; adjusted filePath for original video is %#+v", event.ID, event.OriginalVideo.AdjustedFilePath)

	stat, err := os.Stat(event.OriginalVideo.AdjustedFilePath)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return fmt.Errorf("expected %#+v to be a file but it was a folder", event.OriginalVideo.AdjustedFilePath)
	}

	log.Printf("event.ID: %v; original video is %v MB", event.ID, float64(stat.Size())/1000000)

	rawData, err := ffmpeg.Probe(event.OriginalVideo.AdjustedFilePath)
	if err != nil {
		return fmt.Errorf("failed to probe %#+v for event.ID: %v: %v", event.OriginalVideo.AdjustedFilePath, event.ID, err)
	}

	data := []byte(rawData)

	var ffProbeOutput FFProbeOutput

	err = json.Unmarshal(data, &ffProbeOutput)
	if err != nil {
		return fmt.Errorf("failed to parse ffprobe output %v for event.ID: %v: %v", string(data), event.ID, err)
	}

	event.OriginalVideo.Width = int64(ffProbeOutput.Streams[0].Width)
	event.OriginalVideo.Height = int64(ffProbeOutput.Streams[0].Height)

	fps, err := strconv.ParseInt(strings.Split(ffProbeOutput.Streams[0].RawFPS, "/")[0], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse left portion of %#+v as int", ffProbeOutput.Streams[0].RawFPS)
	}
	event.OriginalVideo.FPS = float64(fps)

	frames, err := strconv.ParseInt(strings.Split(ffProbeOutput.Streams[0].RawFrames, "/")[0], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse left portion of %#+v as int", ffProbeOutput.Streams[0].RawFrames)
	}
	event.OriginalVideo.Frames = frames

	durationSeconds, err := strconv.ParseFloat(ffProbeOutput.Streams[0].RawDuration, 64)
	if err != nil {
		return fmt.Errorf("failed to parse %#+v as int", ffProbeOutput.Streams[0].RawDuration)
	}
	event.OriginalVideo.Duration = time.Second * time.Duration(durationSeconds)

	log.Printf("event.ID: %v; original video is %v frames of %v * %v @ %v FPS for %v",
		event.ID,
		event.OriginalVideo.Frames,
		event.OriginalVideo.Width,
		event.OriginalVideo.Height,
		event.OriginalVideo.FPS,
		event.OriginalVideo.Duration,
	)

	log.Printf("event.ID: %v; preparing %v detections...", event.ID, len(event.PartialDetections))

	event.PartialDetectionsByFrame = make(map[int64][]*PartialDetection, 0)

	for i, detection := range event.PartialDetections {
		detection := detection

		centroidPoints, err := ParseGeometry(detection.RawCentroid)
		if err != nil {
			return fmt.Errorf("failed to parse geometry for centroid %#+v for event.ID: %v", detection.RawCentroid, event.ID)
		}

		if len(centroidPoints) != 1 {
			return fmt.Errorf("failed to get exactly 1 set of points for centroid %#+v for event.ID: %v", detection.RawCentroid, event.ID)
		}
		detection.Centroid = centroidPoints[0]

		detection.BoundingBox, err = ParseGeometry(detection.RawBoundingBox)
		if err != nil {
			return fmt.Errorf("failed to parse geometry for bounding box %#+v for event.ID: %v", detection.RawBoundingBox, event.ID)
		}

		detection.Height = detection.BoundingBox[2].Y - detection.BoundingBox[0].Y
		detection.Width = detection.BoundingBox[2].X - detection.BoundingBox[0].X
		detection.Area = detection.Width * detection.Height
		detection.AspectRatio = detection.Width / detection.Height

		startUnix := event.OriginalVideo.StartTimestamp.UnixNano()
		endUnix := event.OriginalVideo.EndTimestamp.UnixNano()
		detectionUnix := detection.Timestamp.UnixNano()
		detectionProgress := float64(detectionUnix-startUnix) / float64(endUnix-startUnix)
		detection.Frame = int64(math.Round(detectionProgress*float64(event.OriginalVideo.Frames))) + 1

		event.PartialDetections[i] = detection

		for i := int64(0); i < objectDetectionStrideFrames; i++ {
			frame := detection.Frame + i

			partialDetections, ok := event.PartialDetectionsByFrame[frame]
			if !ok {
				partialDetections = make([]*PartialDetection, 0)
			}
			partialDetections = append(partialDetections, detection)

			event.PartialDetectionsByFrame[frame] = partialDetections
		}
	}

	for i, detection := range event.PartialDetections {
		detection.SameClassPartialDetectionsByFrame = make(map[int64][]*PartialDetection)

		for j, otherDetection := range event.PartialDetections {
			if i == j {
				continue
			}

			if detection.ClassID != otherDetection.ClassID {
				continue
			}

			sameClassPartialDetections, ok := detection.SameClassPartialDetectionsByFrame[otherDetection.Frame]
			if !ok {
				sameClassPartialDetections = make([]*PartialDetection, 0)
			}

			sameClassPartialDetections = append(sameClassPartialDetections, otherDetection)

			detection.SameClassPartialDetectionsByFrame[otherDetection.Frame] = sameClassPartialDetections
		}
	}

	log.Printf("event.ID: %v; ready for object tracking.", event.ID)

	return nil
}

func (o *ObjectTracker) trackEvent(event PartialEvent) error {
	if event.OriginalVideo.AdjustedFilePath == "" ||
		event.OriginalVideo.Duration.Seconds() == 0 ||
		event.OriginalVideo.Width == 0 ||
		event.OriginalVideo.Height == 0 {
		return fmt.Errorf("it looks like %#+v hasn't been enriched yet; cannot track", event)
	}

	reader, writer := io.Pipe()
	defer func() {
		_ = reader.Close()
		_ = writer.Close()
	}()

	go func() {
		<-o.ctx.Done()
		_ = reader.Close()
		_ = writer.Close()
	}()

	errs := make(chan error, 16)

	go func() {
		err := ffmpeg.Input(event.OriginalVideo.AdjustedFilePath).
			Output("pipe:",
				ffmpeg.KwArgs{
					"format":  "rawvideo",
					"pix_fmt": "rgb24",
				}).
			WithOutput(writer).
			// ErrorToStdOut().
			Run()
		if err != nil {
			errs <- err
		}
	}()

	select {
	case <-time.After(time.Second * 5):
	case err := <-errs:
		return fmt.Errorf("failed to open ffmpeg input stream for event.ID: %v: %v", event.ID, err)
	}

	width := int(event.OriginalVideo.Width)
	height := int(event.OriginalVideo.Height)
	frameSize := width * height * 3
	buf := make([]byte, frameSize)

	frame := int64(0)

	boundingBoxHoldFrames := int64(boundingBoxHoldDuration.Seconds() * event.OriginalVideo.FPS)
	centroidHoldFrames := int64(centroidHoldDuration.Seconds() * event.OriginalVideo.FPS)
	objectLookbackFrames := int64(objectLookbackDuration.Seconds() * event.OriginalVideo.FPS)

	boundingBoxHoldFramesColorStep := uint8(255 / boundingBoxHoldFrames)
	centroidHoldFramesColorStep := uint8(255 / centroidHoldFrames)

	charSize := gocv.GetTextSize("a", gocv.FontHersheyPlain, 1.0, 1)

	objectID := 0

	for {
		n, err := io.ReadFull(reader, buf)
		if n != frameSize || (err != nil && err != io.EOF) {
			return fmt.Errorf("failed to read %#+v after %v bytes for event.ID: %v: %v", event.OriginalVideo.AdjustedFilePath, n, event.ID, err)
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

		frame++

		overlay := gocv.NewMat()
		gocv.CvtColor(originalImage, &overlay, gocv.ColorBGRToRGBA)
		originalImage.Close()

		frameDurationSeconds := float64(frame) / event.OriginalVideo.FPS
		frameDuration := time.Nanosecond * time.Duration(frameDurationSeconds*1000000000)
		frameTimestamp := event.OriginalVideo.StartTimestamp.Add(frameDuration)

		thisFrameObjectIDs := make(map[int]struct{})

		// infer the objects
		for i, detection := range event.PartialDetectionsByFrame[frame] {

			otherDetectionByScore := make(map[Score]*PartialDetection)

			for i := int64(0); i < objectLookbackFrames; i++ {
				currentFrame := frame - i
				if currentFrame < 1 || currentFrame == frame {
					continue
				}

				for _, otherDetection := range detection.SameClassPartialDetectionsByFrame[currentFrame] {
					if otherDetection.ObjectID != 0 {
						_, alreadySeenThisFrame := thisFrameObjectIDs[int(otherDetection.ObjectID)]
						if alreadySeenThisFrame {
							continue
						}
					}

					centroidDistance := math.Sqrt(
						math.Pow(
							(otherDetection.Centroid.X/float64(width))-
								(detection.Centroid.X/float64(width)),
							2,
						) +
							math.Pow(
								(otherDetection.Centroid.Y/float64(height))-
									(detection.Centroid.Y/float64(height)),
								2,
							),
					)

					score := Score{
						ID:                     otherDetection.ID,
						FrameDistanceFactor:    1.0 - float64(detection.Frame-otherDetection.Frame)/float64(objectLookbackFrames),
						AreaFactor:             detection.Area / otherDetection.Area,
						AspectRatioFactor:      detection.AspectRatio / detection.AspectRatio,
						CentroidDistanceFactor: 1.0 - centroidDistance,
					}

					if 1.0-score.AreaFactor > areaFactorLimit {
						continue
					}

					if 1.0-score.AspectRatioFactor > aspectRatioFactorLimit {
						continue
					}

					if 1.0-score.CentroidDistanceFactor > centroidDistanceFactorLimit {
						continue
					}

					otherDetectionByScore[score] = otherDetection
				}
			}

			scores := maps.Keys(otherDetectionByScore)

			slices.SortFunc(
				scores,
				func(a Score, b Score) int {
					aFrameDistanceFactor := a.FrameDistanceFactor
					if aFrameDistanceFactor == 0 {
						aFrameDistanceFactor = 1.0
					}

					aAreaFactor := a.AreaFactor
					if aAreaFactor == 0 {
						aAreaFactor = 1.0
					}

					aAspectRatioFactor := a.AspectRatioFactor
					if aAspectRatioFactor == 0 {
						aAspectRatioFactor = 1.0
					}

					aCentroidDistanceFactor := a.CentroidDistanceFactor
					if aCentroidDistanceFactor == 0 {
						aCentroidDistanceFactor = 1.0
					}

					bFrameDistanceFactor := b.FrameDistanceFactor
					if bFrameDistanceFactor == 0 {
						bFrameDistanceFactor = 1.0
					}

					bAreaFactor := b.AreaFactor
					if bAreaFactor == 0 {
						bAreaFactor = 1.0
					}

					bAspectRatioFactor := b.AspectRatioFactor
					if bAspectRatioFactor == 0 {
						bAspectRatioFactor = 1.0
					}

					bCentroidDistanceFactor := b.CentroidDistanceFactor
					if bCentroidDistanceFactor == 0 {
						bCentroidDistanceFactor = 1.0
					}

					aScore := aFrameDistanceFactor * aAreaFactor * aAspectRatioFactor * aFrameDistanceFactor
					bScore := bFrameDistanceFactor * bAreaFactor * bAspectRatioFactor * bFrameDistanceFactor

					if aScore < bScore {
						return -1
					} else if aScore > bScore {
						return 1
					}

					return 0
				},
			)

			if len(scores) == 0 {
				objectID++
				detection.ObjectID = int64(objectID)
				event.PartialDetectionsByFrame[frame][i] = detection
			} else {
				detection.ObjectID = otherDetectionByScore[scores[0]].ObjectID
			}

			thisFrameObjectIDs[objectID] = struct{}{}

			log.Printf("%v (%v); %v (%v) has %v scores", frameTimestamp, frame, detection.ClassName, detection.ObjectID, len(scores))
		}

		// draw the bounding box and trail
		for i := int64(0); i < boundingBoxHoldFrames; i++ {
			currentFrame := frame - i

			for _, detection := range event.PartialDetectionsByFrame[currentFrame] {
				topLeft := detection.BoundingBox[0]
				bottomRight := detection.BoundingBox[2]

				boundingBox := image.Rectangle{
					Min: image.Point{
						X: int(topLeft.X),
						Y: int(topLeft.Y),
					},
					Max: image.Point{
						X: int(bottomRight.X),
						Y: int(bottomRight.Y),
					},
				}

				boundingBoxColor := color.RGBA{
					255 - (boundingBoxHoldFramesColorStep * uint8(i)),
					0,
					0,
					0,
				}

				gocv.Rectangle(&overlay, boundingBox, boundingBoxColor, 1)
			}
		}

		// draw the centroid and trail
		for i := int64(0); i < centroidHoldFrames; i++ {
			currentFrame := frame - i

			for _, detection := range event.PartialDetectionsByFrame[currentFrame] {
				centroid := image.Point{int(detection.Centroid.X), int(detection.Centroid.Y)}

				centroidColor := color.RGBA{
					255 - (centroidHoldFramesColorStep * uint8(i)),
					255 - (centroidHoldFramesColorStep * uint8(i)),
					255 - (centroidHoldFramesColorStep * uint8(i)),
					0,
				}

				gocv.Circle(&overlay, centroid, 4, centroidColor, 2)
			}
		}

		// draw the labels
		for _, detection := range event.PartialDetectionsByFrame[frame] {
			bottomLeft := detection.BoundingBox[3]

			textColor := color.RGBA{
				0,
				0,
				0,
				0,
			}

			text1 := fmt.Sprintf("%v (%v) @ %.2f%%", detection.ClassName, detection.ObjectID, detection.Score*100.0)
			// text2 := fmt.Sprintf("%.2fw x %.2fh", detection.Width, detection.Height)
			// text3 := fmt.Sprintf("%.2f | %.2f", detection.Area, detection.AspectRatio)

			textPoint1 := image.Pt(int(bottomLeft.X), int(bottomLeft.Y)+(charSize.Y*0))
			// textPoint2 := image.Pt(int(bottomLeft.X), int(bottomLeft.Y)+(charSize.Y*2))
			// textPoint3 := image.Pt(int(bottomLeft.X), int(bottomLeft.Y)+(charSize.Y*4))

			gocv.PutText(&overlay, text1, textPoint1, gocv.FontHersheyPlain, 1.0, textColor, 1)
			// gocv.PutText(&overlay, text2, textPoint2, gocv.FontHersheyPlain, 1.0, textColor, 1)
			// gocv.PutText(&overlay, text3, textPoint3, gocv.FontHersheyPlain, 1.0, textColor, 1)
		}

		if o.mats != nil {
			o.mats <- overlay
		}
	}

	return nil
}

func (o *ObjectTracker) HandleEvent(event PartialEvent) error {
	err := o.enrichEvent(&event) // note: mutates event
	if err != nil {
		return nil
	}

	err = o.trackEvent(event)
	if err != nil {
		return nil
	}

	return nil
}
