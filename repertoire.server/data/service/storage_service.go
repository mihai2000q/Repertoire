package service

import (
	"bytes"
	"errors"
	"github.com/patrickmn/go-cache"
	"io"
	"mime/multipart"
	"net/http"
	"repertoire/server/internal"
	"time"

	"github.com/go-resty/resty/v2"
)

type StorageService interface {
	Upload(fileHeader *multipart.FileHeader, filePath string) error
	DeleteFile(filePath internal.FilePath) error
	DeleteDirectory(directoryPath string) error
}

type storageService struct {
	httpClient *resty.Client
	env        internal.Env
	cache      *cache.Cache
}

func NewStorageService(httpClient *resty.Client, env internal.Env, cache *cache.Cache) StorageService {
	return &storageService{
		httpClient: httpClient,
		env:        env,
		cache:      cache,
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

	token, err := s.getAccessToken()
	if err != nil {
		return err
	}

	res, err := s.httpClient.R().
		SetAuthToken(token).
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
	token, err := s.getAccessToken()
	if err != nil {
		return err
	}

	res, err := s.httpClient.R().
		SetAuthToken(token).
		Delete("files/" + string(*filePath.StripURL()))
	if err != nil {
		return err
	}
	if res.StatusCode() != http.StatusOK {
		return errors.New("Storage Service - DeleteFile failed: " + res.String())
	}

	return nil
}

func (s storageService) DeleteDirectory(directoryPath string) error {
	token, err := s.getAccessToken()
	if err != nil {
		return err
	}

	res, err := s.httpClient.R().
		SetAuthToken(token).
		Delete("directories/" + directoryPath)
	if err != nil {
		return err
	}
	if res.StatusCode() != http.StatusOK {
		return errors.New("Storage Service - DeleteDirectory failed: " + res.String())
	}

	return nil
}

func (s storageService) getAccessToken() (string, error) {
	// get from cache
	accessTokenKey := "access_token"
	token, found := s.cache.Get(accessTokenKey)
	if found {
		return token.(string), nil
	}

	// fetch from server
	tokenResult, err := s.fetchToken()
	if err != nil {
		return "", err
	}
	expiresIn, _ := time.ParseDuration(tokenResult.expiresIn)
	s.cache.Set(accessTokenKey, tokenResult.accessToken, expiresIn)
	return tokenResult.accessToken, nil
}

type tokenResponse struct {
	accessToken string
	expiresIn   string
}

func (s storageService) fetchToken() (tokenResponse, error) {
	var result tokenResponse
	response, err := s.httpClient.R().
		SetFormData(map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     s.env.StorageClientID,
			"client_secret": s.env.StorageClientSecret,
		}).
		SetResult(&result).
		Post(s.env.AuthStorageUrl)

	if err != nil {
		return tokenResponse{}, err
	}
	if response.StatusCode() != http.StatusOK {
		return tokenResponse{}, errors.New("Storage Service - oauth token failed: " + response.String())
	}

	return result, nil
}
