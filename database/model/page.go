package model

type Page struct {
	Offset *int64 `json:"offset,omitempty"`
	Limit  *int64 `json:"limit,omitempty"`
	Total  uint64 `json:"total"`
}
