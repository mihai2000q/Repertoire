package client

import (
	"github.com/go-resty/resty/v2"
	"io"
	"repertoire/server/data/http"
	"repertoire/server/internal"
)

type StorageClient struct {
	resty.Client
}

func NewStorageClient(client http.RestyClient, env internal.Env) StorageClient {
	return StorageClient{
		*client.SetBaseURL(env.StorageUrl),
	}
}

func (client StorageClient) Upload(token string, fileName string, reader io.Reader, filePath string) (*resty.Response, error) {
	return client.R().
		SetAuthToken(token).
		SetFileReader("file", fileName, reader).
		SetFormData(map[string]string{
			"filePath": filePath,
		}).
		Put("upload")
}

func (client StorageClient) DeleteFile(token string, filePath string) (*resty.Response, error) {
	return client.R().
		SetAuthToken(token).
		Delete("files/" + filePath)
}

func (client StorageClient) DeleteDirectories(token string, directoryPaths []string) (*resty.Response, error) {
	var body = struct{ DirectoryPaths []string }{directoryPaths}
	return client.R().
		SetAuthToken(token).
		SetBody(body).
		Put("directories")
}
