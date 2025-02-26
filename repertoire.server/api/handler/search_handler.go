package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/api/server"
	"repertoire/server/api/validation"
	"repertoire/server/domain/service"
)

type SearchHandler struct {
	service service.SearchService
	server.BaseHandler
}

func NewSearchHandler(
	service service.SearchService,
	validator *validation.Validator,
) *SearchHandler {
	return &SearchHandler{
		service: service,
		BaseHandler: server.BaseHandler{
			Validator: validator,
		},
	}
}

func (s SearchHandler) Get(c *gin.Context) {
	var request requests.SearchGetRequest
	err := c.BindQuery(&request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	errorCode := s.Validator.Validate(&request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	token := s.GetTokenFromContext(c)

	result, errorCode := s.service.Get(request, token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, result)
}
