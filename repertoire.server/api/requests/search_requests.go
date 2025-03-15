package requests

import "repertoire/server/internal/enums"

type SearchGetRequest struct {
	Query       string            `form:"query" validate:"required"`
	CurrentPage *int              `form:"currentPage" validate:"required_with=PageSize,omitempty,gt=0"`
	PageSize    *int              `form:"pageSize" validate:"required_with=CurrentPage,omitempty,gt=0"`
	Type        *enums.SearchType `form:"type" validate:"omitempty,search_type_enum"`
	Filter      []string          `form:"filter"`
	Order       []string          `form:"order"`
}
