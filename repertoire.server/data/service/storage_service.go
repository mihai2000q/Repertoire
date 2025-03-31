package service

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"repertoire/server/data/cache"
	"repertoire/server/data/http/auth"
	"repertoire/server/data/http/client"
	"repertoire/server/internal"
	"repertoire/server/internal/wrapper"
	"strings"
	"time"
)

type StorageService interface {
	Upload(fileHeader *multipart.FileHeader, filePath string) error
	DeleteFile(filePath internal.FilePath) *wrapper.ErrorCode
	DeleteDirectories(directoryPaths []string) *wrapper.ErrorCode
}

type storageService struct {
	storageClient client.StorageClient
	authClient    client.AuthClient
	env           internal.Env
	cache         cache.StorageCache
}

func NewStorageService(
	httpClient client.StorageClient,
	authClient client.AuthClient,
	cache cache.StorageCache,
) StorageService {
	return &storageService{
		storageClient: httpClient,
		authClient:    authClient,
		cache:         cache,
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

	userID := s.getUserIDFromPath(filePath)
	storageToken, err := s.getAccessToken(userID)
	if err != nil {
		return err
	}

	res, err := s.storageClient.Upload(storageToken, fileHeader.Filename, bytes.NewReader(buf.Bytes()), filePath)
	if err != nil {
		return err
	}
	if res.StatusCode() != http.StatusOK {
		return errors.New("Storage Service - Upload failed: " + res.String())
	}

	return nil
}

func (s storageService) DeleteFile(filePath internal.FilePath) *wrapper.ErrorCode {
	stringFilePath := string(*filePath.StripURL())
	userID := s.getUserIDFromPath(stringFilePath)
	storageToken, err := s.getAccessToken(userID)
	if err != nil {
		return wrapper.UnauthorizedError(err)
	}

	res, err := s.storageClient.DeleteFile(storageToken, stringFilePath)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if res.StatusCode() == http.StatusNotFound {
		return wrapper.NotFoundError(errors.New("Storage Service - DeleteFile Not Found: " + res.String()))
	}
	if res.StatusCode() != http.StatusOK {
		return wrapper.InternalServerError(errors.New("Storage Service - DeleteFile failed: " + res.String()))
	}

	return nil
}

func (s storageService) DeleteDirectories(directoryPaths []string) *wrapper.ErrorCode {
	storageToken, err := s.getAccessToken(directoryPaths[0])
	if err != nil {
		return wrapper.UnauthorizedError(err)
	}

	res, err := s.storageClient.DeleteDirectories(storageToken, directoryPaths)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if res.StatusCode() == http.StatusNotFound {
		return wrapper.NotFoundError(errors.New("Storage Service - DeleteDirectories Not Found: " + res.String()))
	}
	if res.StatusCode() != http.StatusOK {
		return wrapper.InternalServerError(errors.New("Storage Service - DeleteDirectories failed: " + res.String()))
	}

	return nil
}

func (s storageService) getAccessToken(userID string) (string, error) {
	// get from cache
	accessTokenKey := "access_token#" + userID
	accessToken, found := s.cache.Get(accessTokenKey)
	if found {
		return accessToken.(string), nil
	}

	// fetch from server
	tokenResult, err := s.fetchToken(userID)
	if err != nil {
		return "", err
	}
	expiresIn, _ := time.ParseDuration(tokenResult.ExpiresIn)
	s.cache.Set(accessTokenKey, tokenResult.Token, expiresIn)
	return tokenResult.Token, nil
}

func (s storageService) fetchToken(userID string) (auth.TokenResponse, error) {
	var result auth.TokenResponse
	response, err := s.authClient.StorageToken(userID, &result)

	if err != nil {
		return auth.TokenResponse{}, err
	}
	if response.StatusCode() != http.StatusOK {
		return auth.TokenResponse{}, errors.New("failed to fetch token: " + response.String())
	}

	return result, nil
}

func (s storageService) getUserIDFromPath(path string) string {
	return strings.Split(path, "/")[0]
}
