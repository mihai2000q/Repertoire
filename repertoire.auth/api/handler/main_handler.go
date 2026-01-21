package handler

import (
	"net/http"
	"repertoire/auth/api/requests"
	"repertoire/auth/api/server"
	"repertoire/auth/api/validation"
	"repertoire/auth/domain/service"

	"github.com/gin-gonic/gin"
)

type MainHandler struct {
	service service.MainService
	server.BaseHandler
}

func NewMainHandler(
	service service.MainService,
	validator *validation.Validator,
) *MainHandler {
	return &MainHandler{
		service: service,
		BaseHandler: server.BaseHandler{
			Validator: validator,
		},
	}
}

func (a MainHandler) Refresh(c *gin.Context) {
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

	c.JSON(http.StatusOK, token)
}

func (a MainHandler) SignIn(c *gin.Context) {
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

	c.JSON(http.StatusOK, token)
}
