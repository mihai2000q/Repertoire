package middleware

import (
	"errors"
	"net/http"
	"repertoire/storage/domain/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtService service.JwtService
}

func NewAuthMiddleware(jwtService service.JwtService) AuthMiddleware {
	return AuthMiddleware{jwtService: jwtService}
}

func (a AuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			err := a.jwtService.Authorize(authToken)
			if err != nil {
				_ = c.AbortWithError(http.StatusUnauthorized, errors.New("invalid token"))
				return
			}

			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "you are not authorized",
		})
	}
}
