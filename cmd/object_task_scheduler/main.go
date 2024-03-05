package main

import (
	"flag"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/initialed85/cameranator/pkg/services/object_task_scheduler"
	"github.com/initialed85/cameranator/pkg/utils"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	urlFlag := flag.String("url", "http://localhost:8080/v1/graphql", "")
	timeoutFlag := flag.Duration("timeout", time.Second*300, "")
	amqpFlag := flag.String("amqp", "amqp://guest:guest@localhost:5672", "")
	skipPublishFlag := flag.String("skipPublish", "false", "")

	flag.Parse()

	url := *urlFlag
	timeout := *timeoutFlag
	amqp := *amqpFlag

	skipPublish, err := strconv.ParseBool(strings.ToLower(*skipPublishFlag))
	if err != nil {
		log.Fatalf("failed to parse -skipPublish %v as bool", *skipPublishFlag)
	}

	if url == "" || !(strings.Contains(url, "http://") || strings.Contains(url, "https://")) {
		log.Fatal("invalid -url argument; must be HTTP URL for GraphQL instance")
	}

	if timeout <= time.Duration(0) {
		log.Fatal("invalid -timeout argument; must be > 0s")
	}

	if amqp == "" || !strings.Contains(amqp, "amqp://") {
		log.Fatal("invalid -amqp argument; must be AMQP URI for RabbitMQ instance")
	}

	objectTaskScheduler, err := object_task_scheduler.NewObjectTaskScheduler(url, timeout, amqp, skipPublish)
	if err != nil {
		log.Fatal(err)
	}

	objectTaskScheduler.Start()
	log.Printf("Press Ctrl + C to exit...")
	utils.WaitForCtrlC()
	objectTaskScheduler.Stop()
}
