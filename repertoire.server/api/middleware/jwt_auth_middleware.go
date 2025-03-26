package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"repertoire/server/data/service"
	"strings"
)

type JWTAuthMiddleware struct {
	jwtService service.JwtService
}

func NewJWTAuthMiddleware(jwtService service.JwtService) JWTAuthMiddleware {
	return JWTAuthMiddleware{jwtService: jwtService}
}

func (m JWTAuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			errorCode := m.jwtService.Authorize(authToken)
			if errorCode != nil {
				_ = c.AbortWithError(errorCode.Code, errorCode.Error)
				return
			} else {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "you are not authorized",
		})
	}
}
