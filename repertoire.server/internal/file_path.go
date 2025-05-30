package internal

import (
	"os"
	"strings"
)

type FilePath string

func (f *FilePath) ToFullURL() *FilePath {
	if f == nil {
		return nil
	}
	if strings.Contains(string(*f), f.getFullURL()) {
		return f
	}
	url := FilePath(f.getFullURL() + string(*f))
	return &url
}

func (f *FilePath) StripURL() *FilePath {
	if f == nil {
		return nil
	}
	url := FilePath(strings.Replace(string(*f), f.getFullURL(), "", -1))
	return &url
}

func (f *FilePath) getFullURL() string {
	return os.Getenv("STORAGE_FETCH_URL")
}
