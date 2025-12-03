package models
type Pagination[T any] struct {
	Page int	`json:"page"`
	LastPage int `json:"last_page"`
	Count	int64	`json:"count"`
	Items []T	`json:"items"`
}