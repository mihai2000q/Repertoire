package handler

import (
	"github.com/gin-gonic/gin"
	"repertoire/auth/api/server"
	"repertoire/auth/api/validation"
	"repertoire/auth/domain/service"
)

type StorageHandler struct {
	service service.StorageService
	server.BaseHandler
}

func NewStorageHandler(
	service service.StorageService,
	validator *validation.Validator,
) *StorageHandler {
	return &StorageHandler{
		service: service,
		BaseHandler: server.BaseHandler{
			Validator: validator,
		},
	}
}

func (s StorageHandler) Token(c *gin.Context) {
	grantType := c.PostForm("grant_type")
	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")
	token := s.GetTokenFromContext(c)

	storageToken, expiresIn, errCode := s.service.Token(grantType, clientID, clientSecret, token)
	if errCode != nil {
		_ = c.AbortWithError(errCode.Code, errCode.Error)
		return
	}

	s.SendToken(c, storageToken, expiresIn)
}
