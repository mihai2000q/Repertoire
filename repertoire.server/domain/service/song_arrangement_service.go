package service

import (
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/arrangement"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type SongArrangementService interface {
	BulkUpdate(request requests.BulkUpdateSongArrangementsRequest) *wrapper.ErrorCode
	Create(request requests.CreateSongArrangementRequest) *wrapper.ErrorCode
	Delete(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode
	GetAll(request requests.GetSongArrangementsRequest) ([]model.SongArrangement, *wrapper.ErrorCode)
	Move(request requests.MoveSongArrangementRequest) *wrapper.ErrorCode
	UpdateDefault(request requests.UpdateDefaultSongArrangementRequest) *wrapper.ErrorCode
}

type songArrangementService struct {
	createSongArrangement        arrangement.CreateSongArrangement
	deleteSongArrangement        arrangement.DeleteSongArrangement
	getAllSongArrangements       arrangement.GetAllSongArrangements
	moveSongArrangement          arrangement.MoveSongArrangement
	updateDefaultSongArrangement arrangement.UpdateDefaultSongArrangement
	bulkUpdateSongArrangement    arrangement.BulkUpdateSongArrangements
}

func NewSongArrangementService(
	createSongArrangement arrangement.CreateSongArrangement,
	deleteSongArrangement arrangement.DeleteSongArrangement,
	getAllSongArrangements arrangement.GetAllSongArrangements,
	moveSongArrangement arrangement.MoveSongArrangement,
	updateDefaultSongArrangement arrangement.UpdateDefaultSongArrangement,
	bulkUpdateSongArrangement arrangement.BulkUpdateSongArrangements,
) SongArrangementService {
	return &songArrangementService{
		createSongArrangement:        createSongArrangement,
		deleteSongArrangement:        deleteSongArrangement,
		getAllSongArrangements:       getAllSongArrangements,
		moveSongArrangement:          moveSongArrangement,
		updateDefaultSongArrangement: updateDefaultSongArrangement,
		bulkUpdateSongArrangement:    bulkUpdateSongArrangement,
	}
}

func (s *songArrangementService) BulkUpdate(request requests.BulkUpdateSongArrangementsRequest) *wrapper.ErrorCode {
	return s.bulkUpdateSongArrangement.Handle(request)
}

func (s *songArrangementService) Create(request requests.CreateSongArrangementRequest) *wrapper.ErrorCode {
	return s.createSongArrangement.Handle(request)
}

func (s *songArrangementService) Delete(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	return s.deleteSongArrangement.Handle(id, songID)
}

func (s *songArrangementService) GetAll(request requests.GetSongArrangementsRequest) ([]model.SongArrangement, *wrapper.ErrorCode) {
	return s.getAllSongArrangements.Handle(request)
}

func (s *songArrangementService) Move(request requests.MoveSongArrangementRequest) *wrapper.ErrorCode {
	return s.moveSongArrangement.Handle(request)
}

func (s *songArrangementService) UpdateDefault(request requests.UpdateDefaultSongArrangementRequest) *wrapper.ErrorCode {
	return s.updateDefaultSongArrangement.Handle(request)
}
