package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"repertoire/storage/internal"
)

type StorageHandler struct {
	env internal.Env
}

func NewStorageHandler(
	env internal.Env,
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

func (s StorageHandler) DeleteDirectories(c *gin.Context) {
	var request struct{ DirectoryPaths []string }
	err := c.BindJSON(&request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	for _, directoryPath := range request.DirectoryPaths {
		storagePath := filepath.Join(s.env.UploadDirectory, directoryPath)
		_ = os.RemoveAll(storagePath)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "directories have been deleted successfully!",
	})
}

func (s StorageHandler) DeleteFile(c *gin.Context) {
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

func (s StorageHandler) DeleteDirectory(c *gin.Context) {
	directoryPath := c.Param("directoryPath")

	storagePath := filepath.Join(s.env.UploadDirectory, directoryPath)

	if _, err := os.Stat(storagePath); err != nil {
		directory := filepath.Base(storagePath)
		_ = c.AbortWithError(http.StatusNotFound, errors.New(fmt.Sprintf("directory not found: %s", directory)))
		return
	}

	if err := os.RemoveAll(storagePath); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "directory has been deleted successfully!",
	})
}
