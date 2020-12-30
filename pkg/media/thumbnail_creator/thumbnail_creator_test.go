package thumbnail_creator

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/initialed85/cameranator/pkg/test_utils"
)

func TestCreateThumbnail(t *testing.T) {
	dir, err := test_utils.GetTempDir()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(dir, "some_file.jpg")
	defer func() {
		_ = os.Remove(path)
	}()

	err = GetThumbnail(
		test_utils.TestVideoPath,
		path,
	)
	defer func() {
		_ = os.Remove(path)
	}()
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(path)
	if err != nil {
		log.Fatal("during test:", err)
	}
}
