package model

type Camera struct {
	ID        int64  `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	StreamURL string `json:"stream_url,omitempty"`
}

func NewCamera(
	name string,
	streamURL string,
) Camera {
	return Camera{
		Name:      name,
		StreamURL: streamURL,
	}
}
