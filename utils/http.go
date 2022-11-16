package utils

type (
	CommonRes struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	ErrorRes struct {
		Status    string `json:"status"`
		Message   string `json:"message"`
		ErrorCode string `json:"error_code"`
	}
)

const (
	SUCCESS = "success"
	ERROR   = "error"
)
