package handler

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"repertoire/storage/internal"
	"testing"
)

var uploadTestDirectory = "../../Test_Uploads/"

// Utils

func getGinHandler() *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.Default()

	storageHandler := StorageHandler{
		env: internal.Env{
			UploadDirectory: uploadTestDirectory,
		},
	}

	engine.PUT("/upload", storageHandler.Upload)
	engine.GET("/files/*filePath", storageHandler.Get)
	engine.DELETE("/files/*filePath", storageHandler.Delete)

	return engine
}

func createMultipartBody(fileName, filePath string) (*bytes.Buffer, string) {
	tempFile, _ := os.CreateTemp("", fileName)
	defer func(name string) {
		_ = os.Remove(name)
	}(tempFile.Name())

	_, _ = tempFile.WriteString("This is a test file")
	_ = tempFile.Close()

	var requestBody bytes.Buffer
	multiWriter := multipart.NewWriter(&requestBody)

	fileWriter, _ := multiWriter.CreateFormFile("file", tempFile.Name())

	file, _ := os.Open(tempFile.Name())
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	_, _ = file.WriteTo(fileWriter)

	_ = multiWriter.WriteField("filePath", filePath)
	_ = multiWriter.Close()

	return &requestBody, multiWriter.FormDataContentType()
}

func createFile(filePath string, content string) {
	// create directories
	dir := filepath.Dir(filepath.Join(uploadTestDirectory, filePath))
	_ = os.MkdirAll(dir, os.ModePerm)

	file, _ := os.Create(filepath.Join(uploadTestDirectory, filePath))
	_, _ = file.WriteString(content)
	_ = file.Close()
}

// Tests

func TestStorageHandler_Get_WhenFileIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	filePath := "somewhere/else/test-file.txt"

	// when
	req := httptest.NewRequest(http.MethodGet, "/files/"+filePath, nil)
	w := httptest.NewRecorder()
	getGinHandler().ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestStorageHandler_Get_WhenSuccessful_ShouldSaveFile(t *testing.T) {
	// given
	t.Cleanup(func() {
		_ = os.RemoveAll(uploadTestDirectory)
	})

	filePath := "somewhere/test-file.txt"
	fileContent := "This is a test file"
	createFile(filePath, fileContent)

	// when
	req := httptest.NewRequest(http.MethodGet, "/files/"+filePath, nil)
	w := httptest.NewRecorder()
	getGinHandler().ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, fileContent, w.Body.String())
}

func TestStorageHandler_Upload_WhenFileIsMissing_ShouldReturnBadRequest(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodPut, "/upload", nil)

	// when
	w := httptest.NewRecorder()
	getGinHandler().ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStorageHandler_Upload_WhenFilePathIsMissing_ShouldReturnBadRequest(t *testing.T) {
	// given
	tempFile, _ := os.CreateTemp("", "test-file.txt")
	defer func(name string) {
		_ = os.Remove(name)
	}(tempFile.Name())

	var requestBody bytes.Buffer
	multiWriter := multipart.NewWriter(&requestBody)

	fileWriter, _ := multiWriter.CreateFormFile("file", tempFile.Name())
	file, _ := os.Open(tempFile.Name())
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	_, _ = file.WriteTo(fileWriter)

	_ = multiWriter.Close()

	req := httptest.NewRequest(http.MethodPut, "/upload", &requestBody)

	// when
	w := httptest.NewRecorder()
	getGinHandler().ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStorageHandler_Upload_WhenSuccessful_ShouldSaveFile(t *testing.T) {
	// given
	t.Cleanup(func() {
		_ = os.RemoveAll(uploadTestDirectory)
	})

	oldFileName := "old-test-file.txt"
	filePath := "somewhere/else/test-file.txt"
	body, contentType := createMultipartBody(oldFileName, filePath)

	// when
	req := httptest.NewRequest(http.MethodPut, "/upload", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	getGinHandler().ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	// assert file uploaded
	uploadedFilePath := filepath.Join(uploadTestDirectory, filePath)
	_, err := os.Stat(uploadedFilePath)
	assert.NoError(t, err)
}

func TestStorageHandler_Delete_WhenFileIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	filePath := "somewhere/else/test-file.txt"

	// when
	req := httptest.NewRequest(http.MethodDelete, "/files/"+filePath, nil)
	w := httptest.NewRecorder()
	getGinHandler().ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestStorageHandler_Delete_WhenSuccessful_ShouldSaveFile(t *testing.T) {
	// given
	t.Cleanup(func() {
		_ = os.RemoveAll(uploadTestDirectory)
	})

	filePath := "somewhere/test-file.txt"
	createFile(filePath, "")

	// when
	req := httptest.NewRequest(http.MethodDelete, "/files/"+filePath, nil)
	w := httptest.NewRecorder()
	getGinHandler().ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	// assert file deleted
	_, err := os.Stat(filepath.Join(uploadTestDirectory, filePath))
	assert.Error(t, err)
}
