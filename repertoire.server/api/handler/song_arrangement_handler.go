package handler

import (
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/api/server"
	"repertoire/server/api/validation"
	"repertoire/server/domain/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SongArrangementHandler struct {
	service service.SongArrangementService
	server.BaseHandler
}

func NewSongArrangementHandler(
	service service.SongArrangementService,
	validator *validation.Validator,
) *SongArrangementHandler {
	return &SongArrangementHandler{
		service: service,
		BaseHandler: server.BaseHandler{
			Validator: validator,
		},
	}
}

func (s SongArrangementHandler) GetAll(c *gin.Context) {
	// TODO: When Gin fixes it, replace with BindQuery
	querySongId := c.Query("songId")
	songId, err := uuid.Parse(querySongId)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var request requests.GetSongArrangementsRequest
	request.SongID = songId
	errorCode := s.Validator.Validate(&request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	result, errorCode := s.service.GetAll(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s SongArrangementHandler) Create(c *gin.Context) {
	var request requests.CreateSongArrangementRequest
	errorCode := s.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = s.service.Create(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "song arrangement has been created successfully!")
}

func (s SongArrangementHandler) BulkUpdate(c *gin.Context) {
	var request requests.BulkUpdateSongArrangementsRequest
	errorCode := s.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = s.service.BulkUpdate(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "song arrangement has been updated successfully!")
}

func (s SongArrangementHandler) UpdateDefault(c *gin.Context) {
	var request requests.UpdateDefaultSongArrangementRequest
	errorCode := s.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = s.service.UpdateDefault(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "song arrangement has been updated to default!")
}

func (s SongArrangementHandler) Move(c *gin.Context) {
	var request requests.MoveSongArrangementRequest
	errorCode := s.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = s.service.Move(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "song arrangement has been moved successfully!")
}

func (s SongArrangementHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	songID, err := uuid.Parse(c.Param("songID"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	errorCode := s.service.Delete(id, songID)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "song arrangement has been deleted successfully!")
}
