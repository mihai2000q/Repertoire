package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"repertoire/data/service"
	"strings"
)

type JWTAuthMiddleware struct {
	jwtService service.JwtService
	// logger  logger.Logger
}

func NewJWTAuthMiddleware(
	jwtService service.JwtService,
	// logger lib.Logger,
) JWTAuthMiddleware {
	return JWTAuthMiddleware{
		jwtService: jwtService,
		//logger:  logger,
	}
}

// Handler handles middleware functionality
func (m JWTAuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			err := m.jwtService.Authorize(authToken)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				})
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
