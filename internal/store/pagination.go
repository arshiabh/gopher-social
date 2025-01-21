package store

import (
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Order  string `json:"order" validate:"oneof=asc desc"`
	Search string `json:"search" validate:"max=100"`
}

func (fq *PaginatedFeedQuery) Parse(r *http.Request) error {
	qs := r.URL.Query()
	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return err
		}
		fq.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		r, err := strconv.Atoi(offset)
		if err != nil {
			return err
		}
		fq.Offset = r
	}

	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}
	return nil
}
