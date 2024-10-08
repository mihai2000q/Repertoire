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
	user, err := u.currentUserProvider.Get(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user.ID == uuid.Nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Query("email")

	user, err := u.service.GetByEmail(email)

	if err != nil {
		// u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user.ID == uuid.Nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
