package response

type PaginatedListResponse[T any] struct {
	Data       T     `json:"data"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalCount int64 `json:"total_count"`
	TotalPage  int   `json:"total_page"`
}
