package handler

import (
	"github.com/gin-gonic/gin"
	"repertoire/auth/api/server"
	"repertoire/auth/api/validation"
	"repertoire/auth/domain/service"
)

type CentrifugoHandler struct {
	service service.CentrifugoService
	server.BaseHandler
}

func NewCentrifugoHandler(
	service service.CentrifugoService,
	validator *validation.Validator,
) *CentrifugoHandler {
	return &CentrifugoHandler{
		service: service,
		BaseHandler: server.BaseHandler{
			Validator: validator,
		},
	}
}

func (c CentrifugoHandler) Token(ctx *gin.Context) {
	token := c.GetTokenFromContext(ctx)

	centrifugoToken, expiresIn, errCode := c.service.Token(token)
	if errCode != nil {
		_ = ctx.AbortWithError(errCode.Code, errCode.Error)
		return
	}

	c.SendToken(ctx, centrifugoToken, expiresIn)
}
