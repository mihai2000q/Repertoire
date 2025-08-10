package internal

import (
	"os"
	"repertoire/server/internal"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToFullURL_WhenIsNil_ShouldReturnNil(t *testing.T) {
	// given
	var _uut *internal.FilePath

	// when
	result := _uut.ToFullURL()

	// then
	assert.Nil(t, result)
}

func TestToFullURL_WhenURLIsAlreadyFull_ShouldReturnTheFilePathAsItIs(t *testing.T) {
	// given
	storageUrl := "the_storage_url"
	_ = os.Setenv("STORAGE_FETCH_URL", storageUrl)
	_uut := internal.FilePath(storageUrl + "some_file_path")

	// when
	result := _uut.ToFullURL()

	// then
	assert.Equal(t, _uut, *result)
}

func TestToFullURL_WhenSuccessful_ShouldReturnFilePathPrefixedByStorageUrlAndSuffixedByLastModifiedDate(t *testing.T) {
	// given
	storageUrl := "the_storage_url"
	_ = os.Setenv("STORAGE_FETCH_URL", storageUrl)
	_uut := internal.FilePath("some_file_path")

	// when
	result := _uut.ToFullURL()

	// then
	assert.Equal(t, storageUrl+string(_uut), string(*result))
}

func TestStripURL_WhenIsNil_ShouldReturnNil(t *testing.T) {
	// given
	var _uut *internal.FilePath

	// when
	result := _uut.StripURL()

	// then
	assert.Nil(t, result)
}

func TestStripURL_WhenSuccessful_ShouldReturnTheFilePathWithoutTheStorageUrlOrTrailingLastModifiedAt(t *testing.T) {
	// given
	storageUrl := "the_storage_url"
	_ = os.Setenv("STORAGE_FETCH_URL", storageUrl)
	filePath := "some_file_path"
	_uut := internal.FilePath(storageUrl + filePath)

	// when
	result := _uut.StripURL()

	// then
	assert.Equal(t, filePath, string(*result))
}
