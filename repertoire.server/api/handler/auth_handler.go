package handler

import (
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/api/server"
	"repertoire/server/api/validation"
	"repertoire/server/domain/service"

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

	c.JSON(http.StatusOK, token)
}
