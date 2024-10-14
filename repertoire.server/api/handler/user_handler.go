package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"repertoire/api/server"
	"repertoire/api/validation"
	"repertoire/domain/provider"
	"repertoire/domain/service"
)

type UserHandler struct {
	service             service.UserService
	currentUserProvider provider.CurrentUserProvider
	server.BaseHandler
}

func NewUserHandler(
	service service.UserService,
	currentUserProvider provider.CurrentUserProvider,
	validator *validation.Validator,
) *UserHandler {
	return &UserHandler{
		service:             service,
		currentUserProvider: currentUserProvider,
		BaseHandler: server.BaseHandler{
			Validator: validator,
		},
	}
}

func (u UserHandler) GetCurrentUser(c *gin.Context) {
	token := u.GetTokenFromContext(c)

	user, errorCode := u.currentUserProvider.Get(token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u UserHandler) Get(c *gin.Context) {
	paramId := c.Param("id")

	id, err := uuid.Parse(paramId)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, errorCode := u.service.Get(id)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, user)
}
