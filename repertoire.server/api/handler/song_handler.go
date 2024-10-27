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

type SongHandler struct {
	service service.SongService
	server.BaseHandler
}

func NewSongHandler(
	service service.SongService,
	validator *validation.Validator,
) *SongHandler {
	return &SongHandler{
		service: service,
		BaseHandler: server.BaseHandler{
			Validator: validator,
		},
	}
}

func (s SongHandler) Get(c *gin.Context) {
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

func (s SongHandler) GetAll(c *gin.Context) {
	request := requests.GetSongsRequest{
		CurrentPage: s.IntQueryOrNull(c, "currentPage"),
		PageSize:    s.IntQueryOrNull(c, "pageSize"),
		OrderBy:     c.Query("orderBy"),
	}
	errorCode := s.Validator.Validate(&request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}
	token := s.GetTokenFromContext(c)

	result, errorCode := s.service.GetAll(request, token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s SongHandler) GetGuitarTunings(c *gin.Context) {
	token := s.GetTokenFromContext(c)

	result, errorCode := s.service.GetGuitarTunings(token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s SongHandler) Create(c *gin.Context) {
	var request requests.CreateSongRequest
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

	s.SendMessage(c, "song has been created successfully")
}

func (s SongHandler) Update(c *gin.Context) {
	var request requests.UpdateSongRequest
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

	s.SendMessage(c, "song has been updated successfully")
}

func (s SongHandler) Delete(c *gin.Context) {
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

	s.SendMessage(c, "song has been deleted successfully")
}
