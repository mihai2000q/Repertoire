package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"repertoire/storage/utils"
)

type StorageHandler struct {
	env utils.Env
}

func NewStorageHandler(
	env utils.Env,
) *StorageHandler {
	return &StorageHandler{
		env: env,
	}
}

func (s StorageHandler) Get(c *gin.Context) {
	filePath := c.Param("filePath")

	storagePath := filepath.Join(s.env.UploadDirectory, filePath)
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		filename := filepath.Base(storagePath)
		_ = c.AbortWithError(http.StatusNotFound, errors.New(fmt.Sprintf("file not found: %s", filename)))
		return
	}

	c.File(storagePath)
}

func (s StorageHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	filePath := c.PostForm("filePath")

	if err = c.SaveUploadedFile(file, filepath.Join(s.env.UploadDirectory, filePath)); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "file has been uploaded successfully!",
	})
}

func (s StorageHandler) Delete(c *gin.Context) {
	filePath := c.Param("filePath")

	storagePath := filepath.Join(s.env.UploadDirectory, filePath)
	if err := os.Remove(storagePath); err != nil {
		if os.IsNotExist(err) {
			filename := filepath.Base(storagePath)
			_ = c.AbortWithError(http.StatusNotFound, errors.New(fmt.Sprintf("file not found: %s", filename)))
			return
		}
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "file has been deleted successfully!",
	})
}
