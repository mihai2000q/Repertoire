package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"repertoire/storage/domain/service"
)

type StorageHandler struct {
	service service.StorageService
}

func NewStorageHandler(
	service service.StorageService,
) *StorageHandler {
	return &StorageHandler{
		service: service,
	}
}

func (u StorageHandler) Get(c *gin.Context) {
	c.JSON(http.StatusOK, "it works")
}
