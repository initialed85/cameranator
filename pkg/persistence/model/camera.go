package model

import (
	"github.com/google/uuid"

	"github.com/initialed85/cameranator/pkg/utils"
)

type Camera struct {
	ID        int64     `json:"id,omitempty"`
	UUID      uuid.UUID `json:"uuid,omitempty"`
	Name      string    `json:"name,omitempty"`
	StreamURL string    `json:"stream_url,omitempty"`
}

func NewCamera(
	Name string,
	StreamURL string,
) Camera {
	randomUUID := utils.GetUUID()

	return Camera{
		UUID:      randomUUID,
		Name:      Name,
		StreamURL: StreamURL,
	}
}
