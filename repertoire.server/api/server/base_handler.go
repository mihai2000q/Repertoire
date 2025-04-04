package server

import (
	"net/http"
	"repertoire/server/api/validation"
	"repertoire/server/internal/wrapper"
	"strings"

	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
	Validator *validation.Validator
}

func (*BaseHandler) GetTokenFromContext(c *gin.Context) string {
	authHeader := c.Request.Header.Get("Authorization")
	t := strings.Split(authHeader, " ")
	return t[1]
}

func (h *BaseHandler) BindAndValidate(c *gin.Context, request any) *wrapper.ErrorCode {
	err := c.Bind(&request)
	if err != nil {
		return wrapper.BadRequestError(err)
	}

	errCode := h.Validator.Validate(request)
	if errCode != nil {
		return errCode
	}

	return nil
}

func (*BaseHandler) SendMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
