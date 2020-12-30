package main

import (
	"flag"
	"log"

	"github.com/initialed85/cameranator/pkg/segments/segment_generator"
	"github.com/initialed85/cameranator/pkg/services/segment_generators"
	"github.com/initialed85/cameranator/pkg/utils"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	netCamURLs := utils.FlagSliceString{}
	cameraNames := utils.FlagSliceString{}

	destinationPathFlag := flag.String("destinationPath", "", "")
	durationFlag := flag.Int("duration", 0, "")
	hostFlag := flag.String("host", "localhost", "")
	portFlag := flag.Int64("port", 6291, "")
	flag.Var(&netCamURLs, "netCamURL", "")
	flag.Var(&cameraNames, "cameraName", "")

	flag.Parse()

	destinationPath := *destinationPathFlag
	duration := *durationFlag
	host := *hostFlag
	port := *portFlag

	if destinationPath == "" {
		log.Fatal("invalid -destinationPath argument; may not be empty")
	}

	if duration <= 0 {
		log.Fatal("invalid -duration argument; must be > 0")
	}

	if host == "" {
		log.Fatal("invalid -host argument; may not be empty")
	}

	if port <= 0 {
		log.Fatal("invalid -port argument; must be > 0")
	}

	if len(netCamURLs) == 0 {
		log.Fatal("invalid -netCamURL argument; need at least 1")
	}

	if len(cameraNames) == 0 {
		log.Fatal("invalid -cameraName argument; need at least 1")
	}

	if len(netCamURLs) != len(cameraNames) {
		log.Fatal("invalid -netCamURL and -cameraName arguments; must have same amount of both (they're indexed together)")
	}

	feeds := make([]segment_generator.Feed, 0)

	for i, netCamURL := range netCamURLs {
		cameraName := cameraNames[i]

		feeds = append(feeds, segment_generator.Feed{
			NetCamURL:       netCamURL,
			DestinationPath: destinationPath,
			CameraName:      cameraName,
			Duration:        duration,
		})
	}

	segmentGenerator := segment_generators.NewSegmentGenerators(
		feeds,
		host,
		port,
	)

	err := segmentGenerator.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Press Ctrl + C to exit...")
	utils.WaitForCtrlC()
	segmentGenerator.Stop()
}
