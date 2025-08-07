package middleware

import (
	"errors"
	"net/http"
	"repertoire/server/internal"

	"github.com/gin-gonic/gin"
)

type MeiliAuthMiddleware struct {
	env internal.Env
}

func NewMeiliAuthMiddleware(env internal.Env) MeiliAuthMiddleware {
	return MeiliAuthMiddleware{env: env}
}

func (m MeiliAuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authKey := c.Request.Header.Get("Authorization")
		if authKey != m.env.MeiliAuthKey {
			_ = c.AbortWithError(http.StatusUnauthorized, errors.New("you are unauthorized"))
			return
		} else {
			c.Next()
			return
		}
	}
}
