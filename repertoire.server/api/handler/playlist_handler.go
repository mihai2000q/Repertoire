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

type PlaylistHandler struct {
	service service.PlaylistService
	server.BaseHandler
}

func NewPlaylistHandler(
	service service.PlaylistService,
	validator *validation.Validator,
) *PlaylistHandler {
	return &PlaylistHandler{
		service: service,
		BaseHandler: server.BaseHandler{
			Validator: validator,
		},
	}
}

func (p PlaylistHandler) Get(c *gin.Context) {
	var request requests.GetPlaylistRequest
	err := c.BindQuery(&request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	request.ID = id

	errorCode := p.Validator.Validate(&request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	playlist, errorCode := p.service.Get(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, playlist)
}

func (p PlaylistHandler) GetAll(c *gin.Context) {
	var request requests.GetPlaylistsRequest
	err := c.BindQuery(&request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	errorCode := p.Validator.Validate(&request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	token := p.GetTokenFromContext(c)

	result, errorCode := p.service.GetAll(request, token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (p PlaylistHandler) GetFiltersMetadata(c *gin.Context) {
	var request requests.GetPlaylistFiltersMetadataRequest
	err := c.BindQuery(&request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	errorCode := p.Validator.Validate(&request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	token := p.GetTokenFromContext(c)

	result, errorCode := p.service.GetFiltersMetadata(request, token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (p PlaylistHandler) Create(c *gin.Context) {
	var request requests.CreatePlaylistRequest
	errorCode := p.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	token := p.GetTokenFromContext(c)

	id, errorCode := p.service.Create(request, token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (p PlaylistHandler) AddAlbums(c *gin.Context) {
	var request requests.AddAlbumsToPlaylistRequest
	errorCode := p.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	res, errorCode := p.service.AddAlbums(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (p PlaylistHandler) AddArtists(c *gin.Context) {
	var request requests.AddArtistsToPlaylistRequest
	errorCode := p.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	res, errorCode := p.service.AddArtists(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (p PlaylistHandler) AddSongs(c *gin.Context) {
	var request requests.AddSongsToPlaylistRequest
	errorCode := p.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	res, errorCode := p.service.AddSongs(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (p PlaylistHandler) Update(c *gin.Context) {
	var request requests.UpdatePlaylistRequest
	errorCode := p.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = p.service.Update(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	p.SendMessage(c, "playlist has been updated successfully")
}

func (p PlaylistHandler) MoveSong(c *gin.Context) {
	var request requests.MoveSongFromPlaylistRequest
	errorCode := p.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = p.service.MoveSong(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	p.SendMessage(c, "song has been moved from playlist successfully")
}

func (p PlaylistHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	errorCode := p.service.Delete(id)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	p.SendMessage(c, "playlist has been deleted successfully")
}

func (p PlaylistHandler) RemoveSongs(c *gin.Context) {
	var request requests.RemoveSongsFromPlaylistRequest
	errorCode := p.BindAndValidate(c, &request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	errorCode = p.service.RemoveSongs(request)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	p.SendMessage(c, "songs have been removed from playlist successfully")
}

// Images

func (p PlaylistHandler) SaveImage(c *gin.Context) {
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

	errorCode := p.service.SaveImage(file, id)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	p.SendMessage(c, "image has been saved to playlist successfully!")
}

func (p PlaylistHandler) DeleteImage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	errorCode := p.service.DeleteImage(id)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	p.SendMessage(c, "image has been deleted from playlist successfully")
}
