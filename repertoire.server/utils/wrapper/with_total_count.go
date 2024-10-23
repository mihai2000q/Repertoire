package wrapper

type WithTotalCount[T any] struct {
	Models     []T   `json:"models"`
	TotalCount int64 `json:"totalCount"`
}
