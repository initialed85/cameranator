package page_renderer

import (
	"log"
	"testing"
	"time"

	"github.com/relvacode/iso8601"
)

func TestNewPageRenderer(t *testing.T) {
	// TODO
}

func TestTimezones(t *testing.T) {
	timestamp := time.Now()

	isoTimestamp := iso8601.Time{Time: timestamp}

	loc, _ := time.LoadLocation("UTC")
	utcTimestamp := timestamp.In(loc)
	utcIsoTimestamp := isoTimestamp.In(loc)

	log.Printf("%v", timestamp)
	log.Printf("%v", isoTimestamp)
	log.Printf("%v", utcTimestamp)
	log.Printf("%v", utcIsoTimestamp)
}
