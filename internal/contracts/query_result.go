package contracts

type QueryResult[T any] struct {
	TotalCount int `json:"totalCount"`
	Items []T `json:"items"`
}
