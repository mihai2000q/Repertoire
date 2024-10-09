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

func (a AuthHandler) Refresh(c *gin.Context) {
	var request auth.RefreshRequest
	err := c.Bind(&request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	token, errCode := a.service.Refresh(request)
	if errCode != nil {
		_ = c.AbortWithError(errCode.Code, errCode.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (a AuthHandler) SignIn(c *gin.Context) {
	var request auth.SignInRequest
	err := c.Bind(&request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	token, errCode := a.service.SignIn(request)
	if errCode != nil {
		_ = c.AbortWithError(errCode.Code, errCode.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (a AuthHandler) SignUp(c *gin.Context) {
	var request auth.SignUpRequest
	err := c.Bind(&request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	token, errCode := a.service.SignUp(request)
	if errCode != nil {
		_ = c.AbortWithError(errCode.Code, errCode.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
