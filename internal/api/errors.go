package api

func NewErrorResponse(code ErrorResponseErrorCode, message string) ErrorResponse {
	return ErrorResponse{
		Error: struct {
			Code    ErrorResponseErrorCode `json:"code"`
			Message string                 `json:"message"`
		}{
			Code:    code,
			Message: message,
		},
	}
}
