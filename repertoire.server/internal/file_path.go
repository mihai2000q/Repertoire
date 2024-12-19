package internal

import (
	"os"
	"strings"
	"time"
)

type FilePath string

func (f *FilePath) ToFullURL(lastModifiedAt *time.Time) *FilePath {
	if f == nil {
		return nil
	}
	url := FilePath(f.getFullURL() + string(*f) + "?lastModifiedAt=" + lastModifiedAt.String())
	return &url
}

func (f *FilePath) StripURL() *FilePath {
	if f == nil {
		return nil
	}
	url := FilePath(strings.Replace(string(*f), f.getFullURL(), "", -1))
	url = FilePath(strings.Split(string(url), "?")[0])
	return &url
}

func (f *FilePath) getFullURL() string {
	return os.Getenv("FETCH_STORAGE_URL") + "/files/"
}
