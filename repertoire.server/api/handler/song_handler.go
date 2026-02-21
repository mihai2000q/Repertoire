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

	song, errorCode := s.service.Get(id)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, song)
}

func (s SongHandler) GetAll(c *gin.Context) {
	var request requests.GetSongsRequest
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

	result, errorCode := s.service.GetAll(request, token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s SongHandler) GetFiltersMetadata(c *gin.Context) {
	var request requests.GetSongFiltersMetadataRequest
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

	result, errorCode := s.service.GetFiltersMetadata(request, token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s SongHandler) Create(c *gin.Context) {
	var request requests.CreateSongRequest
	errCode := s.BindAndValidate(c, &request)
	if errCode != nil {
		_ = c.AbortWithError(errCode.Code, errCode.Error)
		return
	}

	token := s.GetTokenFromContext(c)

	id, errCode := s.service.Create(request, token)
	if errCode != nil {
		_ = c.AbortWithError(errCode.Code, errCode.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (s SongHandler) AddPerfectRehearsal(c *gin.Context) {
	var request requests.AddPerfectSongRehearsalRequest
	errorCode := s.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = s.service.AddPerfectRehearsal(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "perfect rehearsal has been added successfully!")
}

func (s SongHandler) AddPerfectRehearsals(c *gin.Context) {
	var request requests.AddPerfectSongRehearsalsRequest
	errorCode := s.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = s.service.AddPerfectRehearsals(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "perfect rehearsals have been added successfully!")
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

func (s SongHandler) UpdateSettings(c *gin.Context) {
	var request requests.UpdateSongSettingsRequest
	errorCode := s.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = s.service.UpdateSettings(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "song's settings have been updated successfully")
}

func (s SongHandler) BulkDelete(c *gin.Context) {
	var request requests.BulkDeleteSongsRequest
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

	s.SendMessage(c, "songs have been deleted successfully")
}

func (s SongHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
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

// Images

func (s SongHandler) SaveImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id, err := uuid.Parse(c.PostForm("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	errorCode := s.service.SaveImage(file, id)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "image has been saved to song successfully!")
}

func (s SongHandler) DeleteImage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	errorCode := s.service.DeleteImage(id)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	s.SendMessage(c, "image has been deleted from song successfully")
}

// Guitar Tunings

func (s SongHandler) GetGuitarTunings(c *gin.Context) {
	token := s.GetTokenFromContext(c)

	result, errorCode := s.service.GetGuitarTunings(token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, result)
}

// Instruments

func (s SongHandler) GetInstruments(c *gin.Context) {
	token := s.GetTokenFromContext(c)

	result, errorCode := s.service.GetInstruments(token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, result)
}
