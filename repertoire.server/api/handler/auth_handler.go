package handler

import (
	"net/http"
	"repertoire/api/requests"
	"repertoire/api/server"
	"repertoire/api/validation"
	"repertoire/domain/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
	server.BaseHandler
}

func NewAuthHandler(service service.AuthService, validator *validation.Validator) *AuthHandler {
	return &AuthHandler{
		service: service,
		BaseHandler: server.BaseHandler{
			Validator: validator,
		},
	}
}

func (a AuthHandler) Refresh(c *gin.Context) {
	var request requests.RefreshRequest
	errCode := a.BindAndValidate(c, &request)
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
	var request requests.SignInRequest
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
	var request requests.SignUpRequest
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
