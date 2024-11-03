package song

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddSongToAlbum_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongToAlbum{repository: songRepository}

	request := requests.AddSongToAlbumRequest{
		ID:      uuid.New(),
		AlbumID: uuid.New(),
	}

	// given - mocking
	song := &model.Song{ID: request.ID}
	songRepository.On("Get", mock.IsType(song), request.ID).
		Return(nil, song).
		Once()

	var count *int64 = &[]int64{12}[0]
	songRepository.On("CountByAlbum", mock.IsType(count), &request.AlbumID).
		Return(nil, count).
		Once()

	songRepository.On("Update", mock.IsType(song)).
		Run(func(args mock.Arguments) {
			newSong := args.Get(0).(*model.Song)
			assert.Equal(t, *newSong.AlbumID, request.AlbumID)
			assert.Equal(t, *newSong.AlbumTrackNo, uint(*count)+1)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
}
