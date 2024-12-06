package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"repertoire/storage/domain/service"
	"repertoire/storage/internal"
)

type AuthHandler struct {
	env        internal.Env
	jwtService service.JwtService
}

func NewAuthHandler(
	env internal.Env,
	jwtService service.JwtService,
) *AuthHandler {
	return &AuthHandler{
		env:        env,
		jwtService: jwtService,
	}
}

func (a AuthHandler) Token(c *gin.Context) {
	grantType := c.PostForm("grant_type")
	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")

	if grantType != "client_credentials" || clientID != a.env.ClientID || clientSecret != a.env.ClientSecret {
		_ = c.AbortWithError(http.StatusUnauthorized, errors.New("you are not authorized"))
		return
	}

	token, err := a.jwtService.CreateToken()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": token,
		"tokenType":   "Bearer",
		"expiresIn":   a.env.JwtExpirationTime,
	})
}
