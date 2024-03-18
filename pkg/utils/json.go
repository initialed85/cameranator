package utils

import (
	"encoding/json"
)

func UnsafeJSONPrettyFormat(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")

	return string(b)
}

func UnsafeJSONFormat(v any) string {
	b, _ := json.Marshal(v)

	return string(b)
}
