package album

import (
	"encoding/json"
	"errors"
	"repertoire/server/domain/message/handler/album"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAlbumsUpdatedHandler_WhenGetAlbumFails_ShouldReturnError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewAlbumsUpdatedHandler(albumRepository, nil)

	ids := []uuid.UUID{uuid.New()}

	internalError := errors.New("internal error")
	albumRepository.On("GetAllByIDsWithSongsAndArtist", new([]model.Album), ids).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(ids)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, internalError, err)

	albumRepository.AssertExpectations(t)
}

func TestAlbumsUpdatedHandler_WhenPublishFails_ShouldReturnError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := album.NewAlbumsUpdatedHandler(albumRepository, messagePublisherService)

	ids := []uuid.UUID{uuid.New()}

	mockAlbums := []model.Album{{ID: uuid.New()}}
	albumRepository.On("GetAllByIDsWithSongsAndArtist", new([]model.Album), ids).
		Return(nil, &mockAlbums).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.UpdateFromSearchEngineTopic, mock.IsType([]any{})).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(ids)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, internalError, err)

	albumRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestAlbumsUpdatedHandler_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := album.NewAlbumsUpdatedHandler(albumRepository, messagePublisherService)

	ids := []uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.New(),
		uuid.New(),
	}

	albums := []model.Album{
		{
			ID: ids[0],
			Songs: []model.Song{
				{ID: uuid.New()},
				{ID: uuid.New()},
				{ID: uuid.New()},
			},
		},
		{ID: ids[1]},
		{
			ID:    ids[2],
			Songs: []model.Song{{ID: uuid.New()}},
		},
		{ID: ids[3]},
	}

	albumRepository.On("GetAllByIDsWithSongsAndArtist", new([]model.Album), ids).
		Return(nil, &albums).
		Once()

	messagePublisherService.On("Publish", topics.UpdateFromSearchEngineTopic, mock.IsType([]any{})).
		Run(func(args mock.Arguments) {
			searches := args.Get(1).([]any)

			length := 0
			for _, a := range albums {
				assert.Contains(t, searches[length].(model.AlbumSearch).ID, a.ID.String())
				length++
				for _, s := range a.Songs {
					assert.Contains(t, searches[length].(model.SongSearch).ID, s.ID.String())
					assert.Equal(t, searches[length].(model.SongSearch).Album.ID, a.ID)
					length++
				}
			}
			assert.Len(t, searches, length)
		}).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(ids)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	albumRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}
