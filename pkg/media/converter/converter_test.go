package converter

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/initialed85/cameranator/pkg/test_utils"
)

func TestConvertVideo(t *testing.T) {
	dir, err := test_utils.GetTempDir()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(dir, "some_file.mp4")
	defer func() {
		_ = os.Remove(path)
	}()

	stdout, stderr, err := ConvertVideo(
		test_utils.TestVideoPath,
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
	dir, err := test_utils.GetTempDir()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(dir, "some_file.jpg")
	defer func() {
		_ = os.Remove(path)
	}()

	stdout, stderr, err := ConvertImage(
		test_utils.TestImagePath,
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

func testNewConverter(
	t *testing.T,
	newConverter func(int, int) *Converter,
	fileName string,
	sourcePath string,
) {
	dir, err := test_utils.GetTempDir()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(dir, fileName)
	defer func() {
		_ = os.Remove(path)
	}()

	results := make([]struct {
		Work Work
		Err  error
	}, 0)

	c := newConverter(
		4,
		16,
	)

	c.Start()
	defer c.Stop()

	time.Sleep(time.Millisecond * 100)

	work := Work{
		SourcePath:      sourcePath,
		DestinationPath: path,
		Width:           640,
		Height:          360,
	}

	c.Submit(
		work,
		func(work Work, err error) {
			result := struct {
				Work Work
				Err  error
			}{
				work,
				err,
			}

			results = append(
				results,
				result,
			)
		},
	)

	timeout := time.Now().Add(time.Second * 10)
	for len(results) < 1 && time.Now().Before(timeout) {
		time.Sleep(time.Millisecond * 100)
	}

	if len(results) < 1 {
		log.Fatal("results empty")
	}

	result := results[0]
	if result.Err != nil {
		log.Fatal(result.Err)
	}

	assert.Equal(t, work, result.Work)

	_, err = os.Stat(path)
	if err != nil {
		log.Fatal("during test:", err)
	}
}

func TestNewVideoConverter(t *testing.T) {
	testNewConverter(t, NewVideoConverter, "some_file.mp4", test_utils.TestVideoPath)
}

func TestNewImageConverter(t *testing.T) {
	testNewConverter(t, NewImageConverter, "some_file.jpg", test_utils.TestImagePath)
}
