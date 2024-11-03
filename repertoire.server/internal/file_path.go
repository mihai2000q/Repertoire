package internal

import (
	"os"
)

type FilePath string

func (i *FilePath) ToFullURL() FilePath {
	return FilePath(os.Getenv("FETCH_STORAGE_URL") + "/files/" + string(*i))
}

func (i *FilePath) ToNullableFullURL() *FilePath {
	if i == nil {
		return nil
	}
	url := FilePath(os.Getenv("FETCH_STORAGE_URL") + "/files/" + string(*i))
	return &url
}
