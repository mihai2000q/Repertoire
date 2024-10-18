package wrapper

type WithTotalCount[T any] struct {
	Data       []T
	TotalCount int64
}
