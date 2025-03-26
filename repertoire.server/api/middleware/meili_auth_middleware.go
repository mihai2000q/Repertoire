package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"repertoire/server/internal"
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
