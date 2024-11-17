package internal

import (
	"os"
	"strings"
)

type FilePath string

func (i *FilePath) ToFullURL() FilePath {
	return FilePath(i.getFullURL() + string(*i))
}

func (i *FilePath) ToNullableFullURL() *FilePath {
	if i == nil {
		return nil
	}
	url := FilePath(i.getFullURL() + string(*i))
	return &url
}

func (i *FilePath) StripNullableURL() *FilePath {
	if i == nil {
		return nil
	}
	url := FilePath(strings.Replace(string(*i), i.getFullURL(), "", -1))
	return &url
}

func (i *FilePath) getFullURL() string {
	return os.Getenv("FETCH_STORAGE_URL") + "/files/"
}
