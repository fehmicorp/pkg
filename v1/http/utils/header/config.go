package header

import "encoding/json"

type RouteConfig struct {
	Active    bool           `json:"isEnabled"`
	Type      string         `json:"type"`
	Auth      Authentication `json:"authentication,omitempty"`
	RateLimit RateLimit      `json:"rateLimit,omitempty"`
	Cors      CORSConfig     `json:"cors,omitempty"`
}

type Request struct {
	Header Header   `json:"header,omitempty"`
	Params []Params `json:"params,omitempty"`
}

type Params struct {
	Key       string `json:"key"`
	Type      string `json:"type,omitempty"`
	Mandatory bool   `json:"required,omitempty"`
	Pattern   string `json:"Pattern,omitempty"`
}

type Header struct {
	Required []string `json:"required"`
	Defaults struct {
		ContentType string `json:"Content-Type"`
	} `json:"defaults,omitempty"`
}

type Authentication struct {
	Active      bool     `json:"enabled"`
	Type        string   `json:"type"`
	Location    string   `json:"location"`
	KeyName     string   `json:"keyName"`
	AllowedKeys []string `json:"allowedKeys"`
}

type CORSConfig struct {
	Active           bool            `json:"enabled"`
	Origins          json.RawMessage `json:"origins"`
	Methods          []string        `json:"methods"`
	Headers          json.RawMessage `json:"headers"`
	AllowCredentials bool            `json:"allowCredentials"`
	MaxAge           int             `json:"maxAge"`
}

type RateLimit struct {
	Active   bool   `json:"enabled"`
	Requests int    `json:"requests"`
	Duration string `json:"duration"`
	Burst    int    `json:"burst"`
	KeyBy    string `json:"keyBy"`
}
