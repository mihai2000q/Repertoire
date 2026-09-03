package service

import (
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/part"
	"repertoire/server/internal/httperror"

	"github.com/google/uuid"
)

type SongPartService interface {
	BulkUpdate(request requests.BulkUpdateSongPartsRequest) *httperror.ErrorCode
	BulkDelete(request requests.BulkDeleteSongPartsRequest) *httperror.ErrorCode
	Create(request requests.CreateSongPartRequest) *httperror.ErrorCode
	Delete(id uuid.UUID, songID uuid.UUID) *httperror.ErrorCode
	MoveInSong(request requests.MoveSongPartInSongRequest) *httperror.ErrorCode
	UpdateAll(request requests.UpdateAllSongPartsRequest) *httperror.ErrorCode
	Update(request requests.UpdateSongPartRequest) *httperror.ErrorCode
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

func (s *songPartService) BulkUpdate(request requests.BulkUpdateSongPartsRequest) *httperror.ErrorCode {
	return s.bulkRehearsalsSongParts.Handle(request)
}

func (s *songPartService) BulkDelete(request requests.BulkDeleteSongPartsRequest) *httperror.ErrorCode {
	return s.bulkDeleteSongParts.Handle(request)
}

func (s *songPartService) Create(request requests.CreateSongPartRequest) *httperror.ErrorCode {
	return s.createSongPart.Handle(request)
}

func (s *songPartService) Delete(id uuid.UUID, songID uuid.UUID) *httperror.ErrorCode {
	return s.deleteSongPart.Handle(id, songID)
}

func (s *songPartService) MoveInSong(request requests.MoveSongPartInSongRequest) *httperror.ErrorCode {
	return s.moveSongPartInSong.Handle(request)
}

func (s *songPartService) UpdateAll(request requests.UpdateAllSongPartsRequest) *httperror.ErrorCode {
	return s.updateAllSongParts.Handle(request)
}

func (s *songPartService) Update(request requests.UpdateSongPartRequest) *httperror.ErrorCode {
	return s.updateSongPart.Handle(request)
}
