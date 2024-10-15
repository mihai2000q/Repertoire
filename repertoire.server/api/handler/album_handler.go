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

type AlbumHandler struct {
	service service.AlbumService
	server.BaseHandler
}

func NewAlbumHandler(
	service service.AlbumService,
	validator *validation.Validator,
) *AlbumHandler {
	return &AlbumHandler{
		service: service,
		BaseHandler: server.BaseHandler{
			Validator: validator,
		},
	}
}

func (s AlbumHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, errorCode := s.service.Get(id)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s AlbumHandler) GetAll(c *gin.Context) {
	userId, err := uuid.Parse(c.Query("userId"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	request := requests.GetAlbumsRequest{
		UserID: userId,
	}
	errorCode := s.Validator.Validate(&request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	albums, errorCode := s.service.GetAll(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, albums)
}

func (s AlbumHandler) Create(c *gin.Context) {
	var request requests.CreateAlbumRequest
	errorCode := s.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	token := s.GetTokenFromContext(c)

	errorCode = s.service.Create(request, token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "album has been created successfully")
}

func (s AlbumHandler) Update(c *gin.Context) {
	var request requests.UpdateAlbumRequest
	errorCode := s.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = s.service.Update(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "album has been updated successfully")
}

func (s AlbumHandler) Delete(c *gin.Context) {
	paramId := c.Param("id")

	id, err := uuid.Parse(paramId)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	errorCode := s.service.Delete(id)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "album has been deleted successfully")
}
