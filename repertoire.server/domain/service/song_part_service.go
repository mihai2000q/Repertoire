package service

import (
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/part"
	"repertoire/server/internal/wrapper"

	"github.com/google/uuid"
)

type SongPartService interface {
	BulkUpdate(request requests.BulkUpdateSongPartsRequest) *wrapper.ErrorCode
	BulkDelete(request requests.BulkDeleteSongPartsRequest) *wrapper.ErrorCode
	Create(request requests.CreateSongPartRequest) *wrapper.ErrorCode
	Delete(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode
	MoveInSong(request requests.MoveSongPartInSongRequest) *wrapper.ErrorCode
	UpdateAll(request requests.UpdateAllSongPartsRequest) *wrapper.ErrorCode
	Update(request requests.UpdateSongPartRequest) *wrapper.ErrorCode
}

type songPartService struct {
	bulkRehearsalsSongParts part.BulkUpdateSongParts
	bulkDeleteSongParts     part.BulkDeleteSongParts
	createSongPart          part.CreateSongPart
	deleteSongPart          part.DeleteSongPart
	moveSongPartInSong      part.MoveSongPartInSong
	updateAllSongParts      part.UpdateAllSongParts
	updateSongPart          part.UpdateSongPart
}

func NewSongPartService(
	bulkRehearsalsSongParts part.BulkUpdateSongParts,
	bulkDeleteSongParts part.BulkDeleteSongParts,
	createSongPart part.CreateSongPart,
	deleteSongPart part.DeleteSongPart,
	moveSongPartInSong part.MoveSongPartInSong,
	updateAllSongParts part.UpdateAllSongParts,
	updateSongPart part.UpdateSongPart,
) SongPartService {
	return &songPartService{
		bulkRehearsalsSongParts: bulkRehearsalsSongParts,
		bulkDeleteSongParts:     bulkDeleteSongParts,
		createSongPart:          createSongPart,
		deleteSongPart:          deleteSongPart,
		moveSongPartInSong:      moveSongPartInSong,
		updateAllSongParts:      updateAllSongParts,
		updateSongPart:          updateSongPart,
	}
}

func (s *songPartService) BulkUpdate(request requests.BulkUpdateSongPartsRequest) *wrapper.ErrorCode {
	return s.bulkRehearsalsSongParts.Handle(request)
}

func (s *songPartService) BulkDelete(request requests.BulkDeleteSongPartsRequest) *wrapper.ErrorCode {
	return s.bulkDeleteSongParts.Handle(request)
}

func (s *songPartService) Create(request requests.CreateSongPartRequest) *wrapper.ErrorCode {
	return s.createSongPart.Handle(request)
}

func (s *songPartService) Delete(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	return s.deleteSongPart.Handle(id, songID)
}

func (s *songPartService) MoveInSong(request requests.MoveSongPartInSongRequest) *wrapper.ErrorCode {
	return s.moveSongPartInSong.Handle(request)
}

func (s *songPartService) UpdateAll(request requests.UpdateAllSongPartsRequest) *wrapper.ErrorCode {
	return s.updateAllSongParts.Handle(request)
}

func (s *songPartService) Update(request requests.UpdateSongPartRequest) *wrapper.ErrorCode {
	return s.updateSongPart.Handle(request)
}
