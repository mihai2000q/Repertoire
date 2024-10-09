package utils

import "net/http"

type ErrorCode struct {
	Error error
	Code  int
}

func BadRequestError(err error) *ErrorCode {
	return &ErrorCode{
		Error: err,
		Code:  http.StatusBadRequest,
	}
}

func UnauthorizedError(err error) *ErrorCode {
	return &ErrorCode{
		Error: err,
		Code:  http.StatusUnauthorized,
	}
}

func NotFoundError(err error) *ErrorCode {
	return &ErrorCode{
		Error: err,
		Code:  http.StatusNotFound,
	}
}

func InternalServerError(err error) *ErrorCode {
	return &ErrorCode{
		Error: err,
		Code:  http.StatusInternalServerError,
	}
}
