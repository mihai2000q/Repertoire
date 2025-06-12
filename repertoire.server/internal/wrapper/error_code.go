package wrapper

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

func ForbiddenError(err error) *ErrorCode {
	return &ErrorCode{
		Error: err,
		Code:  http.StatusForbidden,
	}
}

func NotFoundError(err error) *ErrorCode {
	return &ErrorCode{
		Error: err,
		Code:  http.StatusNotFound,
	}
}

func ConflictError(err error) *ErrorCode {
	return &ErrorCode{
		Error: err,
		Code:  http.StatusConflict,
	}
}

func InternalServerError(err error) *ErrorCode {
	return &ErrorCode{
		Error: err,
		Code:  http.StatusInternalServerError,
	}
}
