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

type SongSectionHandler struct {
	service service.SongSectionService
	server.BaseHandler
}

func NewSongSectionHandler(
	service service.SongSectionService,
	validator *validation.Validator,
) *SongSectionHandler {
	return &SongSectionHandler{
		service: service,
		BaseHandler: server.BaseHandler{
			Validator: validator,
		},
	}
}

func (s SongSectionHandler) Create(c *gin.Context) {
	var request requests.CreateSongSectionRequest
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

	s.SendMessage(c, "song section has been created successfully!")
}

func (s SongSectionHandler) Update(c *gin.Context) {
	var request requests.UpdateSongSectionRequest
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

	s.SendMessage(c, "song section has been updated successfully!")
}

func (s SongSectionHandler) UpdateOccurrences(c *gin.Context) {
	var request requests.UpdateSongSectionsOccurrencesRequest
	errorCode := s.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = s.service.UpdateOccurrences(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "song sections' occurrences have been updated successfully!")
}

func (s SongSectionHandler) UpdatePartialOccurrences(c *gin.Context) {
	var request requests.UpdateSongSectionsPartialOccurrencesRequest
	errorCode := s.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = s.service.UpdatePartialOccurrences(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "song sections' partial occurrences have been updated successfully!")
}

func (s SongSectionHandler) UpdateAll(c *gin.Context) {
	var request requests.UpdateAllSongSectionsRequest
	errorCode := s.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = s.service.UpdateAll(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "song's sections have been updated successfully based on settings!")
}

func (s SongSectionHandler) Move(c *gin.Context) {
	var request requests.MoveSongSectionRequest
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

	s.SendMessage(c, "song section has been moved successfully!")
}

func (s SongSectionHandler) BulkDelete(c *gin.Context) {
	var request requests.BulkDeleteSongSectionsRequest
	errorCode := s.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = s.service.BulkDelete(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "song sections has been deleted successfully!")
}

func (s SongSectionHandler) Delete(c *gin.Context) {
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

	s.SendMessage(c, "song section has been deleted successfully!")
}

// Types

func (s SongSectionHandler) GetTypes(c *gin.Context) {
	token := s.GetTokenFromContext(c)

	result, errorCode := s.service.GetTypes(token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, result)
}
