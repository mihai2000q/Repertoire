package handler

import (
	"github.com/gin-gonic/gin"
	"repertoire/auth/api/server"
	"repertoire/auth/api/validation"
	"repertoire/auth/domain/service"
	"repertoire/auth/model"
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
	clientCredentials := model.ClientCredentials{
		GrantType:    c.PostForm("grant_type"),
		ClientID:     c.PostForm("client_id"),
		ClientSecret: c.PostForm("client_secret"),
	}
	token := s.GetTokenFromContext(c)

	storageToken, expiresIn, errCode := s.service.Token(clientCredentials, token)
	if errCode != nil {
		_ = c.AbortWithError(errCode.Code, errCode.Error)
		return
	}

	s.SendToken(c, storageToken, expiresIn)
}
