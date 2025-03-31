package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"repertoire/auth/data/service"
	"strings"
)

type AuthMiddleware struct {
	jwtService service.JwtService
}

func NewAuthMiddleware(jwtService service.JwtService) AuthMiddleware {
	return AuthMiddleware{
		jwtService: jwtService,
	}
}

func (m AuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			errCode := m.jwtService.Authorize(authToken)
			if errCode != nil {
				_ = c.AbortWithError(errCode.Code, errCode.Error)
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
