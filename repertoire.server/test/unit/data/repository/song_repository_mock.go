package repository

import (
	"repertoire/server/model"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
)

type SongRepositoryMock struct {
	mock.Mock
}

func (s *SongRepositoryMock) Get(song *model.Song, id uuid.UUID) error {
	args := s.Called(song, id)

	if len(args) > 1 {
		*song = *args.Get(1).(*model.Song)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) GetWithPlaylistsAndSongs(song *model.Song, id uuid.UUID) error {
	args := s.Called(song, id)

	if len(args) > 1 {
		*song = *args.Get(1).(*model.Song)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) GetWithSections(song *model.Song, id uuid.UUID) error {
	args := s.Called(song, id)

	if len(args) > 1 {
		*song = *args.Get(1).(*model.Song)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) GetWithAssociations(song *model.Song, id uuid.UUID) error {
	args := s.Called(song, id)

	if len(args) > 1 {
		*song = *args.Get(1).(*model.Song)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) GetAllByAlbumAndTrackNo(songs *[]model.Song, albumID uuid.UUID, trackNo uint) error {
	args := s.Called(songs, albumID, trackNo)

	if len(args) > 1 {
		*songs = *args.Get(1).(*[]model.Song)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) GetAllByIDs(songs *[]model.Song, ids []uuid.UUID) error {
	args := s.Called(songs, ids)

	if len(args) > 1 {
		*songs = *args.Get(1).(*[]model.Song)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) GetAllByIDsWithSongs(songs *[]model.Song, ids []uuid.UUID) error {
	args := s.Called(songs, ids)

	if len(args) > 1 {
		*songs = *args.Get(1).(*[]model.Song)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) GetAllByUser(
	songs *[]model.Song,
	userID uuid.UUID,
	currentPage *int,
	pageSize *int,
	orderBy []string,
	searchBy []string,
) error {
	args := s.Called(songs, userID, currentPage, pageSize, orderBy, searchBy)

	if len(args) > 1 {
		*songs = *args.Get(1).(*[]model.Song)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error {
	args := s.Called(count, userID, searchBy)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) CountByAlbum(count *int64, albumID *uuid.UUID) error {
	args := s.Called(count, albumID)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) Create(song *model.Song) error {
	args := s.Called(song)
	return args.Error(0)
}

func (s *SongRepositoryMock) Update(song *model.Song) error {
	args := s.Called(song)
	return args.Error(0)
}

func (s *SongRepositoryMock) UpdateLastTimePlayed(songID uuid.UUID, lastTimePlayed time.Time) error {
	args := s.Called(songID, lastTimePlayed)
	return args.Error(0)
}

func (s *SongRepositoryMock) UpdateAll(songs *[]model.Song) error {
	args := s.Called(songs)
	return args.Error(0)
}

func (s *SongRepositoryMock) UpdateWithAssociations(song *model.Song) error {
	args := s.Called(song)
	return args.Error(0)
}

func (s *SongRepositoryMock) UpdateAllWithAssociations(songs *[]model.Song) error {
	args := s.Called(songs)
	return args.Error(0)
}

func (s *SongRepositoryMock) Delete(id uuid.UUID) error {
	args := s.Called(id)
	return args.Error(0)
}

// Guitar Tunings

func (s *SongRepositoryMock) GetGuitarTunings(tunings *[]model.GuitarTuning, userID uuid.UUID) error {
	args := s.Called(tunings, userID)

	if len(args) > 1 {
		*tunings = *args.Get(1).(*[]model.GuitarTuning)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) GetGuitarTuningsCount(count *int64, userID uuid.UUID) error {
	args := s.Called(count, userID)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) CreateGuitarTuning(tuning *model.GuitarTuning) error {
	args := s.Called(tuning)
	return args.Error(0)
}

func (s *SongRepositoryMock) UpdateAllGuitarTunings(tunings *[]model.GuitarTuning) error {
	args := s.Called(tunings)
	return args.Error(0)
}

func (s *SongRepositoryMock) DeleteGuitarTuning(id uuid.UUID) error {
	args := s.Called(id)
	return args.Error(0)
}

// Sections

func (s *SongRepositoryMock) GetSection(section *model.SongSection, id uuid.UUID) error {
	args := s.Called(section, id)

	if len(args) > 1 {
		*section = *args.Get(1).(*model.SongSection)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) CountSectionsBySong(count *int64, songID uuid.UUID) error {
	args := s.Called(count, songID)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) CreateSection(section *model.SongSection) error {
	args := s.Called(section)
	return args.Error(0)
}

func (s *SongRepositoryMock) UpdateSection(section *model.SongSection) error {
	args := s.Called(section)
	return args.Error(0)
}

func (s *SongRepositoryMock) DeleteSection(id uuid.UUID) error {
	args := s.Called(id)
	return args.Error(0)
}

// Section Types

func (s *SongRepositoryMock) GetSectionTypes(tunings *[]model.SongSectionType, userID uuid.UUID) error {
	args := s.Called(tunings, userID)

	if len(args) > 1 {
		*tunings = *args.Get(1).(*[]model.SongSectionType)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) CountSectionTypes(count *int64, userID uuid.UUID) error {
	args := s.Called(count, userID)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) CreateSectionType(sectionType *model.SongSectionType) error {
	args := s.Called(sectionType)
	return args.Error(0)
}

func (s *SongRepositoryMock) UpdateAllSectionTypes(sectionTypes *[]model.SongSectionType) error {
	args := s.Called(sectionTypes)
	return args.Error(0)
}

func (s *SongRepositoryMock) DeleteSectionType(id uuid.UUID) error {
	args := s.Called(id)
	return args.Error(0)
}

// History

func (s *SongRepositoryMock) GetSongSectionHistory(
	history *[]model.SongSectionHistory,
	sectionID uuid.UUID,
	property model.SongSectionProperty,
) error {
	args := s.Called(history, sectionID, property)

	if len(args) > 1 {
		*history = *args.Get(1).(*[]model.SongSectionHistory)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) CreateSongSectionHistory(history *model.SongSectionHistory) error {
	args := s.Called(history)
	return args.Error(0)
}
