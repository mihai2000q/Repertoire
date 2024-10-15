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

type ArtistHandler struct {
	service service.ArtistService
	server.BaseHandler
}

func NewArtistHandler(
	service service.ArtistService,
	validator *validation.Validator,
) *ArtistHandler {
	return &ArtistHandler{
		service: service,
		BaseHandler: server.BaseHandler{
			Validator: validator,
		},
	}
}

func (s ArtistHandler) Get(c *gin.Context) {
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

func (s ArtistHandler) GetAll(c *gin.Context) {
	userId, err := uuid.Parse(c.Query("userId"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	request := requests.GetArtistsRequest{
		UserID: userId,
	}
	errorCode := s.Validator.Validate(&request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	artists, errorCode := s.service.GetAll(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, artists)
}

func (s ArtistHandler) Create(c *gin.Context) {
	var request requests.CreateArtistRequest
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

	s.SendMessage(c, "artist has been created successfully")
}

func (s ArtistHandler) Update(c *gin.Context) {
	var request requests.UpdateArtistRequest
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

	s.SendMessage(c, "artist has been updated successfully")
}

func (s ArtistHandler) Delete(c *gin.Context) {
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

	s.SendMessage(c, "artist has been deleted successfully")
}
