package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/initialed85/cameranator/pkg/services/motion_processor"
	"github.com/initialed85/cameranator/pkg/utils"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	portFlag := flag.Int64("port", 6291, "")
	urlFlag := flag.String("url", "http://localhost:8080/v1/graphql", "")
	timeoutFlag := flag.Duration("timeout", time.Second*30, "")

	flag.Parse()

	port := *portFlag
	url := *urlFlag
	timeout := *timeoutFlag

	if port <= 0 {
		log.Fatal("invalid -port argument; must be > 0")
	}

	if url == "" || !strings.Contains(url, "http://") {
		log.Fatal("invalid -url argument; must be HTTP URL for GraphQL instance")
	}

	if timeout <= time.Duration(0) {
		log.Fatal("invalid -timeout argument; must be > 0s")
	}

	motionProcessor, err := motion_processor.NewMotionProcessor(port, url, timeout)
	if err != nil {
		log.Fatal(err)
	}

	err = motionProcessor.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Press Ctrl + C to exit...")
	utils.WaitForCtrlC()
	motionProcessor.Stop()
}
