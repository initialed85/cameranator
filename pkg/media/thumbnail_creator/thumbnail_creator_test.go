package thumbnail_creator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/initialed85/cameranator/pkg/test_utils"
	"github.com/stretchr/testify/require"
)

func TestCreateThumbnail(t *testing.T) {
	dir, err := test_utils.GetTempDir()
	if err != nil {
		require.NoError(t, err)
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
		require.NoError(t, err)
	}

	_, err = os.Stat(path)
	if err != nil {
		require.NoError(t, err)
	}
}
