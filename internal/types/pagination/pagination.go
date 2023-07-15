package pagination

type Pagination struct {
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
	Count  uint64 `json:"count"`
}
