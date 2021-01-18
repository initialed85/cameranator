package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/initialed85/cameranator/pkg/services/event_pruner"
	"github.com/initialed85/cameranator/pkg/utils"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	urlFlag := flag.String("url", "http://localhost:8080/v1/graphql", "")
	timeoutFlag := flag.Duration("timeout", time.Second*30, "")
	intervalFlag := flag.Duration("interval", time.Minute*5, "")

	flag.Parse()

	url := *urlFlag
	timeout := *timeoutFlag
	interval := *intervalFlag

	if url == "" || !strings.Contains(url, "http://") {
		log.Fatal("invalid -url argument; must be HTTP URL for GraphQL instance")
	}

	if timeout <= time.Duration(0) {
		log.Fatal("invalid -timeout argument; must be > 0s")
	}

	if interval <= time.Duration(0) {
		log.Fatal("invalid -interval argument; must be > 0s")
	}

	eventPruner, err := event_pruner.NewEventPruner(url, timeout, interval)
	if err != nil {
		log.Fatal(err)
	}

	eventPruner.Start()
	log.Printf("Press Ctrl + C to exit...")
	utils.WaitForCtrlC()
	eventPruner.Stop()
}
