package frame

import (
	"github.com/initialed85/cameranator/pkg/persistence/model"
	"github.com/relvacode/iso8601"
)

type Frame struct {
	model.Camera
	Timestamp iso8601.Time `json:"timestamp"`
	Data      []byte
}
