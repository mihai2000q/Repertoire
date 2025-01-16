package server

import (
	"net/http"
	"repertoire/server/api/validation"
	"repertoire/server/internal/wrapper"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BaseHandler struct {
	Validator *validation.Validator
}

func (*BaseHandler) GetTokenFromContext(c *gin.Context) string {
	authHeader := c.Request.Header.Get("Authorization")
	t := strings.Split(authHeader, " ")
	return t[1]
}

func (h *BaseHandler) BindAndValidate(c *gin.Context, request interface{}) *wrapper.ErrorCode {
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

func (h *BaseHandler) UuidQuery(c *gin.Context, str string) uuid.UUID {
	result, err := uuid.Parse(c.Query(str))
	if err != nil {
		result = uuid.Nil
	}
	return result
}

func (h *BaseHandler) IntQueryOrNull(c *gin.Context, str string) *int {
	result, err := strconv.Atoi(c.Query(str))
	if err != nil {
		return nil
	}
	return &result
}

func (*BaseHandler) SendMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
