package dto

type Meta struct {
	Total      int64 `json:"total"`
	PerPage    int64 `json:"per_page"`
	Page       int64 `json:"page"`
	TotalPages int64 `json:"total_pages"`
}
