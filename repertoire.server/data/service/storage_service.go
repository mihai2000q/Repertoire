package service

import (
	"bytes"
	"github.com/go-resty/resty/v2"
	"io"
	"mime/multipart"
)

type StorageService interface {
	Upload(token string, fileHeader *multipart.FileHeader, filePath string) error
}

type storageService struct {
	httpClient *resty.Client
}

func NewStorageService(httpClient *resty.Client) StorageService {
	return &storageService{
		httpClient: httpClient,
	}
}

func (s storageService) Upload(token string, fileHeader *multipart.FileHeader, filePath string) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	_ = file.Close()

	// Read the file content into a buffer
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	if err != nil {
		return err
	}

	_, err = s.httpClient.R().
		SetAuthToken(token).
		SetFileReader("file", fileHeader.Filename, bytes.NewReader(buf.Bytes())).
		SetFormData(map[string]string{
			"filePath": filePath,
		}).
		Put("upload")

	return err
}

func (s storageService) Delete(token string, filePath string) error {
	_, err := s.httpClient.R().
		SetAuthToken(token).
		Delete("files" + filePath)

	return err
}
