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

func (a AlbumHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, errorCode := a.service.Get(id)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (a AlbumHandler) GetAll(c *gin.Context) {
	request := requests.GetAlbumsRequest{
		CurrentPage: a.IntQueryOrNull(c, "currentPage"),
		PageSize:    a.IntQueryOrNull(c, "pageSize"),
		OrderBy:     c.Query("orderBy"),
	}
	errorCode := a.Validator.Validate(&request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}
	token := a.GetTokenFromContext(c)

	result, errorCode := a.service.GetAll(request, token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (a AlbumHandler) Create(c *gin.Context) {
	var request requests.CreateAlbumRequest
	errorCode := a.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	token := a.GetTokenFromContext(c)

	errorCode = a.service.Create(request, token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	a.SendMessage(c, "album has been created successfully")
}

func (a AlbumHandler) AddSong(c *gin.Context) {
	var request requests.AddSongToAlbumRequest
	errorCode := a.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = a.service.AddSong(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	a.SendMessage(c, "song has been added to album successfully")
}

func (a AlbumHandler) Update(c *gin.Context) {
	var request requests.UpdateAlbumRequest
	errorCode := a.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = a.service.Update(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	a.SendMessage(c, "album has been updated successfully")
}

func (a AlbumHandler) MoveSong(c *gin.Context) {
	var request requests.MoveSongFromAlbumRequest
	errorCode := a.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = a.service.MoveSong(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	a.SendMessage(c, "song has been moved from album successfully")
}

func (a AlbumHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	errorCode := a.service.Delete(id)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	a.SendMessage(c, "album has been deleted successfully")
}

func (a AlbumHandler) RemoveSong(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	songID, err := uuid.Parse(c.Param("songId"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	errorCode := a.service.RemoveSong(id, songID)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	a.SendMessage(c, "song has been removed from album successfully")
}

// Images

func (a AlbumHandler) SaveImage(c *gin.Context) {
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

	errorCode := a.service.SaveImage(file, id)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	a.SendMessage(c, "image has been saved to album successfully!")
}

func (a AlbumHandler) DeleteImage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	errorCode := a.service.DeleteImage(id)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	a.SendMessage(c, "image has been deleted from album successfully")
}
