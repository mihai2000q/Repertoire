package storage

import (
	"github.com/stretchr/testify/assert"
	"repertoire/server/internal/message/topics"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestDeleteDirectoriesStorage_WhenSuccessful_ShouldDeleteDirectory(t *testing.T) {
	// given
	directories := []string{"some_directory", "some_other_directory"}

	// when
	err := utils.PublishToTopic(topics.DeleteDirectoriesStorageTopic, directories)

	// then
	assert.NoError(t, err)
}
