package storage

import (
	"repertoire/server/internal/message/topics"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteDirectoriesStorage_WhenSuccessful_ShouldDeleteDirectory(t *testing.T) {
	// given
	directories := []string{"some_directory", "some_other_directory"}

	// when
	err := utils.PublishToTopic(topics.DeleteDirectoriesStorageTopic, directories)

	// then
	assert.NoError(t, err)
}
