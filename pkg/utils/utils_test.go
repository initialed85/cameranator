package utils

import (
	"syscall"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/relvacode/iso8601"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUUID(t *testing.T) {
	assert.NotEqual(
		t,
		uuid.UUID{},
		GetUUID(),
	)
}

func TestGetISO8601Time(t *testing.T) {
	assert.NotEqual(
		t,
		iso8601.Time{},
		GetISO8601Time("2020-03-27T08:30:00+08:00"),
	)
}

func TestWaitForCtrlC(t *testing.T) {
	interrupted := false

	go func() {
		WaitForCtrlC()
		interrupted = true
	}()

	time.Sleep(time.Millisecond * 100)
	assert.False(t, interrupted)

	err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	require.NoError(t, err)

	time.Sleep(time.Millisecond * 100)
	assert.True(t, interrupted)
}
