package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"repertoire/api/contracts/auth"
	"repertoire/domain/service"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (a AuthHandler) SignUp(c *gin.Context) {
	var request auth.SignUpRequest
	err := c.Bind(&request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	token, err := a.service.SignUp(request)
	// TODO: What if the user already exists ? Return Bad Request
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
