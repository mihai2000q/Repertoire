package handler

import (
	"errors"
	"net/http"
	"repertoire/auth/api/server"
	"repertoire/auth/api/validation"
	"repertoire/auth/domain/service"
	"repertoire/auth/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (c CentrifugoHandler) PublicToken(ctx *gin.Context) {
	clientCredentials := model.ClientCredentials{
		GrantType:    ctx.PostForm("grant_type"),
		ClientID:     ctx.PostForm("client_id"),
		ClientSecret: ctx.PostForm("client_secret"),
	}
	userID, err := uuid.Parse(ctx.PostForm("user_id"))
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid user id"))
		return
	}

	centrifugoToken, expiresIn, errCode := c.service.PublicToken(clientCredentials, userID)
	if errCode != nil {
		_ = ctx.AbortWithError(errCode.Code, errCode.Error)
		return
	}

	c.SendToken(ctx, centrifugoToken, expiresIn)
}

func (c CentrifugoHandler) Token(ctx *gin.Context) {
	token := c.GetTokenFromContext(ctx)

	centrifugoToken, _, errCode := c.service.Token(token)
	if errCode != nil {
		_ = ctx.AbortWithError(errCode.Code, errCode.Error)
		return
	}

	ctx.JSON(http.StatusOK, centrifugoToken)
}
