package entity

type PaginatedItems struct {
	TotalItems            int64 `json:"total_items"`
	TotalPages            int64 `json:"total_pages"`
	CurrentPageTotalItems int64 `json:"current_page_total_items"`
	CurrentPage           int64 `json:"current_page"`
	Items                 any    `json:"items"`
}

func NewPaginationInfo(totalItems, totalPages, currentPageTotalItems, currentPage int64, items any) *PaginatedItems {
	return &PaginatedItems{
		totalItems,
		totalPages,
		currentPageTotalItems,
		currentPage,
		items,
	}
}
