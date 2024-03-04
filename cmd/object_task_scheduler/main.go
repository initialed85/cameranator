package main

import (
	"flag"
	"log"
	"strings"

	"github.com/initialed85/cameranator/pkg/services/object_task_scheduler"
	"github.com/initialed85/cameranator/pkg/utils"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	urlFlag := flag.String("url", "http://localhost:8080/v1/graphql", "")
	amqpFlag := flag.String("amqp", "amqp://guest:guest@localhost:5672", "")

	flag.Parse()

	url := *urlFlag
	amqp := *amqpFlag

	if url == "" || !strings.Contains(url, "http://") {
		log.Fatal("invalid -url argument; must be HTTP URL for GraphQL instance")
	}

	if amqp == "" || !strings.Contains(amqp, "amqp://") {
		log.Fatal("invalid -amqp argument; must be AMQP URI for RabbitMQ instance")
	}

	objectTaskScheduler, err := object_task_scheduler.NewObjectTaskScheduler(url, amqp)
	if err != nil {
		log.Fatal(err)
	}

	objectTaskScheduler.Start()
	log.Printf("Press Ctrl + C to exit...")
	utils.WaitForCtrlC()
	objectTaskScheduler.Stop()
}
