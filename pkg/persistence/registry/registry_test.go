package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/initialed85/cameranator/pkg/persistence/model"
)

func TestRegistry_GetModelAndClient(t *testing.T) {
	r := NewRegistry()

	err := r.Register(testGetModel())
	require.NoError(t, err)

	modelAndClient, err := r.GetModelAndClient("camera", testGetClient())
	require.NoError(t, err)

	items := make([]model.Camera, 0)
	err = modelAndClient.GetAll(&items)
	require.NoError(t, err)

	assert.GreaterOrEqual(
		t,
		len(items),
		0,
	)
}
