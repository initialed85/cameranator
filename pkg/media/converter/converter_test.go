package converter

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testVideoPath = "../../../test_data/events/Event_2020-12-27T10:25:05__104__Testing__01.mp4"
	testImagePath = "../../../test_data/events/Event_2020-12-27T10:25:09__104__Testing__01.jpg"
)

func getTempDir() (string, error) {
	return ioutil.TempDir("", "cameranator")
}

func TestConvertVideo(t *testing.T) {
	DisableNvidia()

	dir, err := getTempDir()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(dir, "some_file.mkv")
	defer func() {
		_ = os.Remove(path)
	}()

	stdout, stderr, err := ConvertVideo(
		testVideoPath,
		path,
		640,
		360,
	)
	defer func() {
		_ = os.Remove(path)
	}()
	if err != nil {
		log.Fatalf("stdout=%#+v, stderr=%#+v, err=%#+v", stdout, stderr, err)
	}

	assert.NotEqual(t, "", stderr)

	_, err = os.Stat(path)
	if err != nil {
		log.Fatal("during test:", err)
	}
}

func TestConvertImage(t *testing.T) {
	dir, err := ioutil.TempDir("", "cameranator")
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(dir, "some_file.jpg")
	defer func() {
		_ = os.Remove(path)
	}()

	stdout, stderr, err := ConvertImage(
		testImagePath,
		path,
		640,
		360,
	)
	defer func() {
		_ = os.Remove(path)
	}()
	if err != nil {
		log.Fatalf("during test: %v; stderr= %v", err, stderr)
	}

	assert.Equal(t, "", stderr)
	assert.Equal(t, "", stdout)

	_, err = os.Stat(path)
	if err != nil {
		log.Fatal("during test:", err)
	}
}

func TestNewVideoConverter(t *testing.T) {
	DisableNvidia()

	dir, err := getTempDir()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(dir, "some_file.mkv")
	defer func() {
		_ = os.Remove(path)
	}()

	results := make([]struct {
		Stdout string
		Stderr string
		Err    error
	}, 0)

	workQueue := make(chan Work, 16)

	c := NewVideoConverter(
		workQueue,
		func(work Work, stdout string, stderr string, err error) {
			results = append(
				results,
				struct {
					Stdout string
					Stderr string
					Err    error
				}{
					Stdout: stdout,
					Stderr: stderr,
					Err:    err,
				},
			)
		},
	)

	c.Start()
	defer c.Stop()

	workQueue <- Work{
		SourcePath:      testVideoPath,
		DestinationPath: path,
		Width:           640,
		Height:          360,
	}

	timeout := time.Now().Add(time.Second * 10)
	for len(results) < 1 && time.Now().Before(timeout) {
		time.Sleep(time.Millisecond * 100)
	}

	assert.Len(t, results, 1)

	result := results[0]
	if result.Err != nil {
		if result.Err != nil {
			log.Fatalf("stdout=%#+v, stderr=%#+v, err=%#+v", result.Stdout, result.Stderr, result.Err)
		}
	}

	assert.NotEqual(t, "", result.Stderr)

	_, err = os.Stat(path)
	if err != nil {
		log.Fatal("during test:", err)
	}
}

func TestNewImageConverter(t *testing.T) {
	DisableNvidia()

	dir, err := getTempDir()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(dir, "some_file.jpg")
	defer func() {
		_ = os.Remove(path)
	}()

	results := make([]struct {
		Stdout string
		Stderr string
		Err    error
	}, 0)

	workQueue := make(chan Work, 16)

	c := NewImageConverter(
		workQueue,
		func(work Work, stdout string, stderr string, err error) {
			results = append(
				results,
				struct {
					Stdout string
					Stderr string
					Err    error
				}{
					Stdout: stdout,
					Stderr: stderr,
					Err:    err,
				},
			)
		},
	)

	c.Start()
	defer c.Stop()

	workQueue <- Work{
		SourcePath:      testImagePath,
		DestinationPath: path,
		Width:           640,
		Height:          360,
	}

	timeout := time.Now().Add(time.Second * 10)
	for len(results) < 1 && time.Now().Before(timeout) {
		time.Sleep(time.Millisecond * 100)
	}

	assert.Len(t, results, 1)

	result := results[0]
	if result.Err != nil {
		if result.Err != nil {
			log.Fatalf("stdout=%#+v, stderr=%#+v, err=%#+v", result.Stdout, result.Stderr, result.Err)
		}
	}

	assert.Equal(t, "", result.Stderr)
	assert.Equal(t, "", result.Stdout)

	_, err = os.Stat(path)
	if err != nil {
		log.Fatal("during test:", err)
	}
}
