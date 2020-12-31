package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/initialed85/cameranator/pkg/services/page_renderers"
	"github.com/initialed85/cameranator/pkg/utils"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	urlFlag := flag.String("url", "http://localhost:8080/v1/graphql", "")
	timeoutFlag := flag.Duration("timeout", time.Second*30, "")
	eventsPathFlag := flag.String("eventsPath", "", "")
	segmentsPathFlag := flag.String("segmentsPath", "", "")

	flag.Parse()

	url := *urlFlag
	timeout := *timeoutFlag
	eventsPath := *eventsPathFlag
	segmentsPath := *segmentsPathFlag

	if url == "" || !strings.Contains(url, "http://") {
		log.Fatal("invalid -url argument; must be HTTP URL for GraphQL instance")
	}

	if timeout <= time.Duration(0) {
		log.Fatal("invalid -timeout argument; must be > 0s")
	}

	if eventsPath == "" {
		log.Fatal("invalid -eventsPath argument; may not be empty")
	}

	if segmentsPath == "" {
		log.Fatal("invalid -segmentsPath argument; may not be empty")
	}

	roles := []page_renderers.Role{
		{
			IsSegment: false,
			Path:      eventsPath,
		},
		{
			IsSegment: true,
			Path:      segmentsPath,
		},
	}

	pageRenderers := page_renderers.NewPageRenderers(url, timeout, roles)

	err := pageRenderers.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Press Ctrl + C to exit...")
	utils.WaitForCtrlC()

	pageRenderers.Stop()
}
