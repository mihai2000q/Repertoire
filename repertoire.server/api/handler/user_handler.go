package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"repertoire/api/utils"
	"repertoire/domain/provider"
	"repertoire/domain/service"
)

type UserHandler struct {
	service             service.UserService
	currentUserProvider provider.CurrentUserProvider
}

func NewUserHandler(
	service service.UserService,
	currentUserProvider provider.CurrentUserProvider,
) *UserHandler {
	return &UserHandler{
		service:             service,
		currentUserProvider: currentUserProvider,
	}
}

func (u UserHandler) GetCurrentUser(c *gin.Context) {
	token := utils.GetTokenFromContext(c)

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
