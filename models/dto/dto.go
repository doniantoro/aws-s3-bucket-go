package dto

type ApiResponse struct {
	Code       string      `json:"code,omitempty"`
	Message    string      `json:"messages"`
	Errors     interface{} `json:"errors"`
	Data       interface{} `json:"data"`
	ServerTime string      `json:"server_time"`
}

type ErrorValidation struct {
	Message   string `json:"message"`
	Parameter string `json:"parameter"`
}
