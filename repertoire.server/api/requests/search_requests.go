package requests

import (
	"repertoire/server/internal/enums"
)

type SearchGetRequest struct {
	Query       string            `form:"query"`
	CurrentPage *int              `form:"currentPage" validate:"required_with=PageSize,omitempty,gt=0"`
	PageSize    *int              `form:"pageSize" validate:"required_with=CurrentPage,omitempty,gt=0"`
	Type        *enums.SearchType `form:"type" validate:"required_with=IDs NotIDs,omitempty,search_type_enum"`
	IDs         []string          `form:"ids"`
	NotIDs      []string          `form:"notIds"`
	Filter      []string          `form:"filter" validate:"search_filter"`
	Order       []string          `form:"order" validate:"search_order"`
}
