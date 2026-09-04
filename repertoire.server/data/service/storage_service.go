package service

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"repertoire/server/data/cache"
	"repertoire/server/data/http/auth"
	"repertoire/server/data/http/client"
	"repertoire/server/internal"
	"repertoire/server/internal/env"
	"repertoire/server/internal/httperror"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type StorageService interface {
	Upload(fileHeader *multipart.FileHeader, filePath string) *httperror.ErrorCode
	DeleteFile(filePath internal.FilePath) *httperror.ErrorCode
	DeleteDirectories(directoryPaths []string) *httperror.ErrorCode
}

type storageService struct {
	storageClient client.StorageClient
	authClient    client.AuthClient
	env           env.Env
	cache         cache.StorageCache
}

func NewStorageService(
	storageClient client.StorageClient,
	authClient client.AuthClient,
	cache cache.StorageCache,
) StorageService {
	return &storageService{
		storageClient: storageClient,
		authClient:    authClient,
		cache:         cache,
	}
}

func (s storageService) Upload(fileHeader *multipart.FileHeader, filePath string) *httperror.ErrorCode {
	file, err := fileHeader.Open()
	if err != nil {
		return httperror.InternalServerError(err)
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, file); err != nil {
		return httperror.InternalServerError(err)
	}

	userID := s.getUserIDFromPath(filePath)
	storageToken, err := s.getAccessToken(userID)
	if err != nil {
		return httperror.UnauthorizedError(err)
	}

	res, err := s.storageClient.Upload(storageToken, fileHeader.Filename, bytes.NewReader(buf.Bytes()), filePath)
	return s.checkStorageResponse("upload", res, err)
}

func (s storageService) DeleteFile(filePath internal.FilePath) *httperror.ErrorCode {
	stringFilePath := string(*filePath.StripURL())
	userID := s.getUserIDFromPath(stringFilePath)
	storageToken, err := s.getAccessToken(userID)
	if err != nil {
		return httperror.UnauthorizedError(err)
	}

	res, err := s.storageClient.DeleteFile(storageToken, stringFilePath)
	return s.checkStorageResponse("delete file", res, err)
}

func (s storageService) DeleteDirectories(directoryPaths []string) *httperror.ErrorCode {
	userID := s.getUserIDFromPath(directoryPaths[0])
	storageToken, err := s.getAccessToken(userID)
	if err != nil {
		return httperror.UnauthorizedError(err)
	}

	res, err := s.storageClient.DeleteDirectories(storageToken, directoryPaths)
	return s.checkStorageResponse("delete directories", res, err)
}

func (s storageService) checkStorageResponse(operation string, res *resty.Response, err error) *httperror.ErrorCode {
	if err != nil {
		return httperror.InternalServerError(fmt.Errorf("storage: %s: %w", operation, err))
	}
	switch res.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return httperror.UnauthorizedError(fmt.Errorf("storage: %s: unauthorized: %s", operation, res.String()))
	case http.StatusForbidden:
		return httperror.ForbiddenError(fmt.Errorf("storage: %s: forbidden: %s", operation, res.String()))
	case http.StatusNotFound:
		return httperror.NotFoundError(fmt.Errorf("storage: %s: not found: %s", operation, res.String()))
	default:
		return httperror.InternalServerError(fmt.Errorf("storage: %s: failed: %s", operation, res.String()))
	}
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
		return auth.TokenResponse{}, fmt.Errorf("storage: fetch token: failed: %s", response.String())
	}

	return result, nil
}

func (s storageService) getUserIDFromPath(path string) string {
	return strings.Split(path, "/")[0]
}
