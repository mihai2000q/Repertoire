package server

import (
	"net/http"
	"repertoire/api/validation"
	"repertoire/utils"
	"strconv"
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

func (h *BaseHandler) BindAndValidate(c *gin.Context, request interface{}) *utils.ErrorCode {
	err := c.Bind(&request)
	if err != nil {
		return utils.BadRequestError(err)
	}

	errCode := h.Validator.Validate(request)
	if errCode != nil {
		return errCode
	}

	return nil
}

func (h *BaseHandler) IntQuery(c *gin.Context, str string) int {
	tempString := c.Query(str)
	result, err := strconv.Atoi(tempString)
	if err != nil {
		result = -1
	}
	return result
}

func (*BaseHandler) SendMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
