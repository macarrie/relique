package api_helpers

type PaginationParams struct {
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
}

type PaginatedResponse[P any] struct {
	Pagination PaginationParams `json:"pagination"`

	Count uint64 `json:"count"`
	Data  []P    `json:"data"`
}
