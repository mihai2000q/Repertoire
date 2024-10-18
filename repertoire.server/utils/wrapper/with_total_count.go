package wrapper

type WithTotalCount[T any] struct {
	Models     []T   `json:"model"`
	TotalCount int64 `json:"totalCount"`
}
