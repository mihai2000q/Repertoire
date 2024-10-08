package utils

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func GetTokenFromContext(c *gin.Context) string {
	authHeader := c.Request.Header.Get("Authorization")
	t := strings.Split(authHeader, " ")
	return t[1]
}
