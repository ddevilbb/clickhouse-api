package paginator

import (
	"math"
	"net/url"
	"strconv"
)

const Limit = 10

type Paginator struct {
	Page   uint64 `json:"page"`
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
	Total  uint64 `json:"total"`
	Pages  uint64 `json:"pages"`
}

func NewPaginator(values url.Values, total uint64) *Paginator {
	var (
		err         error
		page, limit uint64
	)

	page, err = strconv.ParseUint(values.Get("page"), 10, 64)
	if err != nil {
		page = 1
	}
	limit, err = strconv.ParseUint(values.Get("limit"), 10, 64)
	if err != nil {
		limit = Limit
	}

	return &Paginator{
		Page:   page,
		Offset: (page - 1) * limit,
		Limit:  limit,
		Total:  total,
		Pages:  uint64(math.Ceil(float64(total) / float64(limit))),
	}
}
