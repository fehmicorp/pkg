package header

import "encoding/json"

type CORSConfig struct {
	Origins json.RawMessage `json:"origins"`
	Methods []string        `json:"methods"`
	Headers json.RawMessage `json:"headers"`
}

type RateLimit struct {
	Requests int    `json:"requests"`
	Duration string `json:"duration"`
	Burst    int    `json:"burst"`
}
