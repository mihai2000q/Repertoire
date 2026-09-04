package httperror

func MessagePublisherError(err error) *ErrorCode {
	return InternalServerError(err)
}
