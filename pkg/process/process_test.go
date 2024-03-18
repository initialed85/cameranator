package process

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunCommand(t *testing.T) {
	stdout, stderr, err := RunCommand("echo", "hello")
	require.NoError(t, err)

	assert.Equal(t, "hello\n", stdout)
	assert.Equal(t, "", stderr)
}

func TestRunBackgroundProcess(t *testing.T) {
	before := time.Now()

	process, err := RunBackgroundProcess("sleep", "1")
	require.NoError(t, err)

	assert.NotNil(t, process)

	assert.NotNil(t, process.Cmd.Process)

	pid1 := process.Cmd.Process.Pid

	time.Sleep(time.Second * 5)

	pid2 := process.Cmd.Process.Pid

	assert.NotEqual(t, pid1, pid2)

	after := time.Now()

	duration := after.Sub(before)

	assert.Greater(t, duration.Seconds(), 5.0)

	process.Stop()
}
