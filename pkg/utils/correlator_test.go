package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCorrelator(t *testing.T) {
	correlator := NewCorrelator()

	complete := false

	correlation := correlator.NewCorrelation(
		func(correlation *Correlation) {
			complete = true
		},
	)

	item1 := correlation.NewItem("Item1")
	item2 := correlation.NewItem("Item2")
	item3 := correlation.NewItem("Item3")

	assert.False(t, complete)

	item1.SetValue(1)
	item1.Complete()

	item2.SetValue("2")
	item2.Complete()

	assert.False(t, complete)

	item3.SetValue(true)
	item3.Complete()

	assert.True(t, complete)
}
