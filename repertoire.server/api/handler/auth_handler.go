package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"repertoire/api/requests/auth"
	"repertoire/api/server"
	"repertoire/api/validation"
	"repertoire/domain/service"
)

type AuthHandler struct {
	service service.AuthService
	server.BaseHandler
}

func NewAuthHandler(service service.AuthService, validator validation.Validator) *AuthHandler {
	return &AuthHandler{
		service: service,
		BaseHandler: server.BaseHandler{
			Validator: validator,
		},
	}
}

func (a AuthHandler) Refresh(c *gin.Context) {
	var request auth.RefreshRequest
	errCode := a.BindAndValidate(c, request)
	if errCode != nil {
		_ = c.AbortWithError(errCode.Code, errCode.Error)
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
	errCode := a.BindAndValidate(c, &request)
	if errCode != nil {
		_ = c.AbortWithError(errCode.Code, errCode.Error)
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
	errCode := a.BindAndValidate(c, &request)
	if errCode != nil {
		_ = c.AbortWithError(errCode.Code, errCode.Error)
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
