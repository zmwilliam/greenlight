package data

import (
	"strings"

	"github.com/zmwilliam/greenlight/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func (f Filters) Validate(v *validator.Validator) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	v.Check(validator.In(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}

func (f Filters) SortValue() string {
	val, _ := strings.CutPrefix(f.Sort, "-")

	for _, v := range f.SortSafelist {
		if val == v {
			return v
		}
	}

	panic("invalid sort value")
}

func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}
