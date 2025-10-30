package handler

import (
	"net/http"
	"repertoire/server/api/requests"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s SongHandler) CreateSection(c *gin.Context) {
	var request requests.CreateSongSectionRequest
	errorCode := s.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = s.service.CreateSection(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "song section has been created successfully!")
}

func (s SongHandler) UpdateSection(c *gin.Context) {
	var request requests.UpdateSongSectionRequest
	errorCode := s.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = s.service.UpdateSection(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "song section has been updated successfully!")
}

func (s SongHandler) UpdateSectionsOccurrences(c *gin.Context) {
	var request requests.UpdateSongSectionsOccurrencesRequest
	errorCode := s.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = s.service.UpdateSectionsOccurrences(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "song sections' occurrences have been updated successfully!")
}

func (s SongHandler) UpdateSectionsPartialOccurrences(c *gin.Context) {
	var request requests.UpdateSongSectionsPartialOccurrencesRequest
	errorCode := s.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = s.service.UpdateSectionsPartialOccurrences(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "song sections' partial occurrences have been updated successfully!")
}

func (s SongHandler) UpdateAllSections(c *gin.Context) {
	var request requests.UpdateAllSongSectionsRequest
	errorCode := s.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = s.service.UpdateAllSections(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "song's sections have been updated successfully based on settings!")
}

func (s SongHandler) MoveSection(c *gin.Context) {
	var request requests.MoveSongSectionRequest
	errorCode := s.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = s.service.MoveSection(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "song section has been moved successfully!")
}

func (s SongHandler) BulkDeleteSections(c *gin.Context) {
	var request requests.BulkDeleteSongSectionsRequest
	errorCode := s.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = s.service.BulkDeleteSections(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "song sections has been deleted successfully!")
}

func (s SongHandler) DeleteSection(c *gin.Context) {
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

	errorCode := s.service.DeleteSection(id, songID)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "song section has been deleted successfully!")
}

// Section Types

func (s SongHandler) GetSectionTypes(c *gin.Context) {
	token := s.GetTokenFromContext(c)

	result, errorCode := s.service.GetSectionTypes(token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, result)
}
