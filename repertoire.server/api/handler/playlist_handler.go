package handler

import (
	"net/http"
	"repertoire/api/requests"
	"repertoire/api/server"
	"repertoire/api/validation"
	"repertoire/domain/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PlaylistHandler struct {
	service service.PlaylistService
	server.BaseHandler
}

func NewPlaylistHandler(
	service service.PlaylistService,
	validator *validation.Validator,
) *PlaylistHandler {
	return &PlaylistHandler{
		service: service,
		BaseHandler: server.BaseHandler{
			Validator: validator,
		},
	}
}

func (p PlaylistHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, errorCode := p.service.Get(id)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (p PlaylistHandler) GetAll(c *gin.Context) {
	request := requests.GetPlaylistsRequest{
		CurrentPage: p.IntQueryOrNull(c, "currentPage"),
		PageSize:    p.IntQueryOrNull(c, "pageSize"),
		OrderBy:     c.Query("orderBy"),
	}
	errorCode := p.Validator.Validate(&request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}
	token := p.GetTokenFromContext(c)

	result, errorCode := p.service.GetAll(request, token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (p PlaylistHandler) Create(c *gin.Context) {
	var request requests.CreatePlaylistRequest
	errorCode := p.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	token := p.GetTokenFromContext(c)

	errorCode = p.service.Create(request, token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	p.SendMessage(c, "playlist has been created successfully")
}

func (p PlaylistHandler) Update(c *gin.Context) {
	var request requests.UpdatePlaylistRequest
	errorCode := p.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = p.service.Update(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	p.SendMessage(c, "playlist has been updated successfully")
}

func (p PlaylistHandler) Delete(c *gin.Context) {
	paramId := c.Param("id")

	id, err := uuid.Parse(paramId)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	errorCode := p.service.Delete(id)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	p.SendMessage(c, "playlist has been deleted successfully")
}
