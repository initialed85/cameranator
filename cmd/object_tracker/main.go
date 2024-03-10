package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/initialed85/cameranator/pkg/services/object_tracker"
	"gocv.io/x/gocv"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	urlFlag := flag.String("url", "http://localhost:8080/v1/graphql", "")
	timeoutFlag := flag.Duration("timeout", time.Second*300, "")

	flag.Parse()

	url := *urlFlag
	timeout := *timeoutFlag

	if url == "" || !(strings.Contains(url, "http://") || strings.Contains(url, "https://")) {
		log.Fatal("invalid -url argument; must be HTTP URL for GraphQL instance")
	}

	if timeout <= time.Duration(0) {
		log.Fatal("invalid -timeout argument; must be > 0s")
	}

	mats := make(chan gocv.Mat)

	objectTaskScheduler, err := object_tracker.NewObjectTracker(url, timeout, mats)
	if err != nil {
		log.Fatal(err)
	}

	window := gocv.NewWindow("cameranator")
	defer func() {
		_ = window.Close()
	}()

	objectTaskScheduler.Start()

	log.Printf("Press Ctrl + C to exit...")

	for {
		mat := <-mats
		window.IMShow(mat)
		mat.Close()
		if window.WaitKey(10) >= 0 {
			break
		}
	}

	objectTaskScheduler.Stop()
}
