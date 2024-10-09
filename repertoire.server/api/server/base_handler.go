package server

import (
	"github.com/gin-gonic/gin"
	"repertoire/api/validation"
	"repertoire/utils"
	"strings"
)

type BaseHandler struct {
	Validator validation.Validator
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
