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

func (a ArtistHandler) Get(c *gin.Context) {
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

func (a ArtistHandler) GetAll(c *gin.Context) {
	request := requests.GetArtistsRequest{
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

func (a ArtistHandler) Create(c *gin.Context) {
	var request requests.CreateArtistRequest
	errorCode := a.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	token := a.GetTokenFromContext(c)

	id, errorCode := a.service.Create(request, token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (a ArtistHandler) AddAlbum(c *gin.Context) {
	var request requests.AddAlbumToArtistRequest
	errorCode := a.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = a.service.AddAlbum(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	a.SendMessage(c, "album has been added to artist successfully")
}

func (a ArtistHandler) AddSong(c *gin.Context) {
	var request requests.AddSongToArtistRequest
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

	a.SendMessage(c, "song has been added to artist successfully")
}

func (a ArtistHandler) Update(c *gin.Context) {
	var request requests.UpdateArtistRequest
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

	a.SendMessage(c, "artist has been updated successfully")
}

func (a ArtistHandler) Delete(c *gin.Context) {
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

	a.SendMessage(c, "artist has been deleted successfully")
}

func (a ArtistHandler) RemoveAlbum(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	albumID, err := uuid.Parse(c.Param("albumID"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	errorCode := a.service.RemoveAlbum(id, albumID)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	a.SendMessage(c, "album has been removed from artist successfully")
}

func (a ArtistHandler) RemoveSong(c *gin.Context) {
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

	errorCode := a.service.RemoveSong(id, songID)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	a.SendMessage(c, "song has been removed from artist successfully")
}

// Images

func (a ArtistHandler) SaveImage(c *gin.Context) {
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

	a.SendMessage(c, "image has been saved to artist successfully!")
}

func (a ArtistHandler) DeleteImage(c *gin.Context) {
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

	a.SendMessage(c, "image has been deleted from artist successfully")
}
