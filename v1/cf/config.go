package cf

type APIResponse[T any] struct {
	Success  bool       `json:"success"`
	Errors   []APIError `json:"errors"`
	Messages []string   `json:"messages"`
	Result   T          `json:"result"`
}

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
