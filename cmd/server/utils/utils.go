package utils

// GenericErrorResponse formats validation errors
type GenericErrorResponse struct {
	Version string            `json:"version"`
	Code    int               `json:"code"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
}

type GenericSuccessResponse struct {
	Version string `json:"version"`
	Code    int    `json:"code"`
}
