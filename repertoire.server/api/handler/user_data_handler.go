package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/api/server"
	"repertoire/server/api/validation"
	"repertoire/server/domain/service"
)

type UserDataHandler struct {
	service service.SongService
	server.BaseHandler
}

func NewUserDataHandler(
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

// Guitar Tunings

func (u UserDataHandler) CreateGuitarTuning(c *gin.Context) {
	var request requests.CreateGuitarTuningRequest
	errorCode := u.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	token := u.GetTokenFromContext(c)

	errorCode = u.service.CreateGuitarTuning(request, token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	u.SendMessage(c, "guitar tuning has been created successfully!")
}

func (u UserDataHandler) MoveGuitarTuning(c *gin.Context) {
	var request requests.MoveGuitarTuningRequest
	errorCode := u.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	token := u.GetTokenFromContext(c)

	errorCode = u.service.MoveGuitarTuning(request, token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	u.SendMessage(c, "guitar tuning has been moved successfully!")
}

func (u UserDataHandler) DeleteGuitarTuning(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	token := u.GetTokenFromContext(c)

	errorCode := u.service.DeleteGuitarTuning(id, token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	u.SendMessage(c, "guitar tuning has been deleted successfully!")
}

// Song Section Types

func (u UserDataHandler) CreateSectionType(c *gin.Context) {
	var request requests.CreateSongSectionTypeRequest
	errorCode := u.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	token := u.GetTokenFromContext(c)

	errorCode = u.service.CreateSectionType(request, token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	u.SendMessage(c, "song section type has been created successfully!")
}

func (u UserDataHandler) MoveSectionType(c *gin.Context) {
	var request requests.MoveSongSectionTypeRequest
	errorCode := u.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	token := u.GetTokenFromContext(c)

	errorCode = u.service.MoveSectionType(request, token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	u.SendMessage(c, "song section type has been moved successfully!")
}

func (u UserDataHandler) DeleteSectionType(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	token := u.GetTokenFromContext(c)

	errorCode := u.service.DeleteSectionType(id, token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	u.SendMessage(c, "song section type has been deleted successfully!")
}
