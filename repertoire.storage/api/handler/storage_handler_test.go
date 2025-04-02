package handler

import (
	"bytes"
	"encoding/json"
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

// Utils

func getGinStorageHandler() (*gin.Engine, internal.Env) {
	gin.SetMode(gin.TestMode)
	engine := gin.Default()

	env := internal.Env{
		UploadDirectory: "../../Test_Uploads/",
	}
	storageHandler := StorageHandler{
		env: env,
	}

	engine.GET("/files/*filePath", storageHandler.Get)
	engine.PUT("/upload", storageHandler.Upload)
	engine.PUT("/directories", storageHandler.DeleteDirectories)
	engine.DELETE("/files/*filePath", storageHandler.DeleteFile)
	engine.DELETE("/directories/*directoryPath", storageHandler.DeleteDirectory)

	return engine, env
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

func createFile(filePath string, content string, uploadTestDirectory string) {
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

	handler, _ := getGinStorageHandler()

	// when
	req := httptest.NewRequest(http.MethodGet, "/files/"+filePath, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestStorageHandler_Get_WhenSuccessful_ShouldReturnFile(t *testing.T) {
	// given
	handler, env := getGinStorageHandler()

	t.Cleanup(func() {
		_ = os.RemoveAll(env.UploadDirectory)
	})

	filePath := "somewhere/test-file.txt"
	fileContent := "This is a test file"
	createFile(filePath, fileContent, env.UploadDirectory)

	// when
	req := httptest.NewRequest(http.MethodGet, "/files/"+filePath, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, fileContent, w.Body.String())
}

func TestStorageHandler_Upload_WhenFileIsMissing_ShouldReturnBadRequest(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodPut, "/upload", nil)

	handler, _ := getGinStorageHandler()

	// when
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

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

	handler, _ := getGinStorageHandler()

	// when
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStorageHandler_Upload_WhenSuccessful_ShouldSaveFile(t *testing.T) {
	// given
	handler, env := getGinStorageHandler()

	t.Cleanup(func() {
		_ = os.RemoveAll(env.UploadDirectory)
	})

	oldFileName := "old-test-file.txt"
	filePath := "somewhere/else/test-file.txt"
	body, contentType := createMultipartBody(oldFileName, filePath)

	// when
	req := httptest.NewRequest(http.MethodPut, "/upload", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	// assert file uploaded
	uploadedFilePath := filepath.Join(env.UploadDirectory, filePath)
	_, err := os.Stat(uploadedFilePath)
	assert.NoError(t, err)
}

func TestStorageHandler_DeleteDirectories_WhenSuccessful_ShouldDeleteDirectory(t *testing.T) {
	// given
	handler, env := getGinStorageHandler()

	t.Cleanup(func() {
		_ = os.RemoveAll(env.UploadDirectory)
	})

	directoryPaths := []string{"somewhere/else", "somewhere/else/two", "somebody/like", "helloThere"}

	for _, directoryPath := range directoryPaths {
		filePath := directoryPath + "/test-file.txt"
		createFile(filePath, "asd", env.UploadDirectory)
	}
	
	directoryPaths = append(directoryPaths, "one/more/directory-that-is-not-to-be-found")

	// when
	var body = struct{ DirectoryPaths []string }{directoryPaths}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/directories", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	// assert directory deleted
	for _, directoryPath := range directoryPaths {
		_, err := os.Stat(filepath.Join(env.UploadDirectory, directoryPath))
		assert.Error(t, err)
	}
}

func TestStorageHandler_DeleteFile_WhenFileIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	filePath := "somewhere/else/test-file.txt"

	handler, _ := getGinStorageHandler()

	// when
	req := httptest.NewRequest(http.MethodDelete, "/files/"+filePath, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestStorageHandler_DeleteFile_WhenSuccessful_ShouldDeleteFile(t *testing.T) {
	// given
	handler, env := getGinStorageHandler()

	t.Cleanup(func() {
		_ = os.RemoveAll(env.UploadDirectory)
	})

	filePath := "somewhere/test-file.txt"
	createFile(filePath, "", env.UploadDirectory)

	// when
	req := httptest.NewRequest(http.MethodDelete, "/files/"+filePath, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	// assert file deleted
	_, err := os.Stat(filepath.Join(env.UploadDirectory, filePath))
	assert.Error(t, err)
}

func TestStorageHandler_DeleteDirectory_WhenDirectoryIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	directoryPath := "somewhere/else"

	handler, _ := getGinStorageHandler()

	// when
	req := httptest.NewRequest(http.MethodDelete, "/directories/"+directoryPath, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestStorageHandler_DeleteDirectory_WhenSuccessful_ShouldDeleteDirectory(t *testing.T) {
	// given
	handler, env := getGinStorageHandler()

	t.Cleanup(func() {
		_ = os.RemoveAll(env.UploadDirectory)
	})

	directoryPath := "somewhere"
	filePath := directoryPath + "/else/test-file.txt"
	createFile(filePath, "", env.UploadDirectory)

	// when
	req := httptest.NewRequest(http.MethodDelete, "/directories/"+directoryPath, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	// assert directory deleted
	_, err := os.Stat(filepath.Join(env.UploadDirectory, directoryPath))
	assert.Error(t, err)
}
