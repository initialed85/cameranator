package filesystem

import (
	"log"
	"regexp"
	"testing"

	"github.com/initialed85/cameranator/pkg/utils"
)

func TestNewWatcher(t *testing.T) {
	matcher, err := regexp.Compile(".*/file\\.txt")
	if err != nil {
		log.Fatal(err)
	}

	w := NewWatcher(
		"/tmp/",
		matcher,
		func(path string, size float64) {
			log.Printf("create %#+v, %#+v", path, size)
		},
		func(path string, size float64) {
			log.Printf("write %#+v, %#+v", path, size)
		},
	)

	w.Start()

	utils.WaitForCtrlC()

	w.Stop()
}
