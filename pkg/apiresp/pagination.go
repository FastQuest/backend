package apiresp

import "math"

type Pagination struct {
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
}

func NewPagination(total int64, currentPage, perPage int) Pagination {
	lastPage := 0
	if perPage > 0 {
		lastPage = int(math.Ceil(float64(total) / float64(perPage)))
	}

	return Pagination{
		Total:       int(total),
		PerPage:     perPage,
		CurrentPage: currentPage,
		LastPage:    lastPage,
	}
}
