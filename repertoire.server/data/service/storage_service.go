package service

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"repertoire/server/internal"

	"github.com/go-resty/resty/v2"
)

type StorageService interface {
	Upload(fileHeader *multipart.FileHeader, filePath string) error
	DeleteFile(filePath internal.FilePath) error
	DeleteDirectory(directoryPath string) error
}

type storageService struct {
	httpClient *resty.Client
}

func NewStorageService(httpClient *resty.Client) StorageService {
	return &storageService{
		httpClient: httpClient,
	}
}

func (s storageService) Upload(fileHeader *multipart.FileHeader, filePath string) error {
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

	res, err := s.httpClient.R().
		SetFileReader("file", fileHeader.Filename, bytes.NewReader(buf.Bytes())).
		SetFormData(map[string]string{
			"filePath": filePath,
		}).
		Put("upload")
	if err != nil {
		return err
	}
	if res.StatusCode() != http.StatusOK {
		return errors.New("Storage Service - Upload failed: " + res.String())
	}

	return nil
}

func (s storageService) DeleteFile(filePath internal.FilePath) error {
	res, err := s.httpClient.R().Delete("files/" + string(filePath.StripURL()))
	if err != nil {
		return err
	}
	if res.StatusCode() != http.StatusOK {
		return errors.New("Storage Service - DeleteFile failed: " + res.String())
	}

	return nil
}

func (s storageService) DeleteDirectory(directoryPath string) error {
	res, err := s.httpClient.R().Delete("directories/" + directoryPath)
	if err != nil {
		return err
	}
	if res.StatusCode() != http.StatusOK {
		return errors.New("Storage Service - DeleteDirectory failed: " + res.String())
	}

	return nil
}
