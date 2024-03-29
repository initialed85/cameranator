package filesystem

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/initialed85/cameranator/pkg/test_utils"
)

func TestNewWatcher(t *testing.T) {
	dir, err := test_utils.GetTempDir()
	require.NoError(t, err)

	matcher, err := regexp.Compile(`.*/file\.txt`)
	require.NoError(t, err)

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
	require.NoError(t, err)
	time.Sleep(time.Millisecond * 100)
	assert.Len(t, created, 1)
	assert.Len(t, wrote, 0)

	_, err = f.Write([]byte("Hello, world."))
	require.NoError(t, err)
	time.Sleep(time.Millisecond * 100)
	assert.Len(t, created, 1)
	assert.Len(t, wrote, 1)

	_, err = f.Write([]byte("Hello, world."))
	require.NoError(t, err)
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
