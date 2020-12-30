package utils

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewExecutor(t *testing.T) {
	jobs := 256

	e := NewExecutor(
		128,
		jobs,
	)
	e.Start()
	defer e.Stop()

	time.Sleep(time.Millisecond * 100)

	values := make([]int, 0)
	mu := sync.Mutex{}

	for i := 0; i < jobs; i++ {
		e.Submit(
			func() (interface{}, error) {
				time.Sleep(time.Second)
				return 1, nil
			},
			func(result interface{}, err error) {
				if err != nil {
					log.Fatal(err)
				}

				mu.Lock()
				values = append(values, result.(int))
				mu.Unlock()
			},
		)
	}

	for len(values) < jobs {
		time.Sleep(time.Millisecond * 100)
	}

	assert.Equal(
		t,
		values[0],
		1,
	)

	assert.Equal(
		t,
		values[jobs-1],
		1,
	)
}
