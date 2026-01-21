package handler

import (
	"errors"
	"net/http"
	"repertoire/auth/api/server"
	"repertoire/auth/api/validation"
	"repertoire/auth/domain/service"
	"repertoire/auth/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	userID, err := uuid.Parse(c.PostForm("user_id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("invalid user id"))
		return
	}

	storageToken, expiresIn, errCode := s.service.Token(clientCredentials, userID)
	if errCode != nil {
		_ = c.AbortWithError(errCode.Code, errCode.Error)
		return
	}

	s.SendToken(c, storageToken, expiresIn)
}
