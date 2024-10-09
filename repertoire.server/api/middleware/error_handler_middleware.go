package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
)

type ErrorHandlerMiddleware struct {
}

func NewErrorHandlerMiddleware() ErrorHandlerMiddleware {
	return ErrorHandlerMiddleware{}
}

func (m ErrorHandlerMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		var errors []string
		for _, err := range c.Errors {
			log.Printf("Error: %+v\n", err)
			errors = append(errors, err.Error())
		}

		if len(errors) == 1 {
			c.JSON(-1, gin.H{
				"error": errors[0],
			})
		} else if len(errors) > 1 {
			c.JSON(-1, gin.H{
				"errors": errors,
			})
		}
	}
}
