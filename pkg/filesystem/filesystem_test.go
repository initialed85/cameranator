package filesystem

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/initialed85/cameranator/pkg/test_utils"
)

func TestNewWatcher(t *testing.T) {
	dir, err := test_utils.GetTempDir()
	if err != nil {
		log.Fatal(err)
	}

	matcher, err := regexp.Compile(".*/file\\.txt")
	if err != nil {
		log.Fatal(err)
	}

	created := make([]File, 0)
	wrote := make([]File, 0)

	w := NewWatcher(
		dir,
		matcher,
		func(file File) {
			log.Printf("created %#+v", file)
			created = append(created, file)
		},
		func(file File) {
			log.Printf("wrote %#+v", file)
			wrote = append(wrote, file)
		},
	)

	w.Start()

	time.Sleep(time.Millisecond * 100)
	assert.Len(t, created, 0)
	assert.Len(t, wrote, 0)

	f, err := os.Create(fmt.Sprintf("%v/file.txt", dir))
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Millisecond * 100)
	assert.Len(t, created, 1)
	assert.Len(t, wrote, 0)

	_, err = f.Write([]byte("Hello, world."))
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Millisecond * 100)
	assert.Len(t, created, 1)
	assert.Len(t, wrote, 1)

	_, err = f.Write([]byte("Hello, world."))
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Millisecond * 100)
	assert.Len(t, created, 1)
	assert.Len(t, wrote, 2)

	assert.Equal(
		t,
		File{
			Name: fmt.Sprintf("%v/file.txt", dir),
			Size: 2.6e-05,
		},
		wrote[1],
	)

	w.Stop()
}
