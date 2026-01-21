package repository

import (
	"repertoire/server/model"

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

func (s *SongRepositoryMock) GetFiltersMetadata(
	metadata *model.SongFiltersMetadata,
	userID uuid.UUID,
	searchBy []string,
) error {
	args := s.Called(metadata, userID, searchBy)

	if len(args) > 1 {
		*metadata = *args.Get(1).(*model.SongFiltersMetadata)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) GetAllByAlbum(songs *[]model.Song, albumID uuid.UUID) error {
	args := s.Called(songs, albumID)

	if len(args) > 1 {
		*songs = *args.Get(1).(*[]model.Song)
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

func (s *SongRepositoryMock) GetAllByIDsWithSections(songs *[]model.Song, ids []uuid.UUID) error {
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

func (s *SongRepositoryMock) GetAllByIDsWithArtistAndAlbum(songs *[]model.Song, ids []uuid.UUID) error {
	args := s.Called(songs, ids)

	if len(args) > 1 {
		*songs = *args.Get(1).(*[]model.Song)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) GetAllByIDsWithAlbumsAndPlaylists(songs *[]model.Song, ids []uuid.UUID) error {
	args := s.Called(songs, ids)

	if len(args) > 1 {
		*songs = *args.Get(1).(*[]model.Song)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) GetAllByUser(
	songs *[]model.EnhancedSong,
	userID uuid.UUID,
	currentPage *int,
	pageSize *int,
	orderBy []string,
	searchBy []string,
) error {
	args := s.Called(songs, userID, currentPage, pageSize, orderBy, searchBy)

	if len(args) > 1 {
		*songs = *args.Get(1).(*[]model.EnhancedSong)
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

func (s *SongRepositoryMock) CountByAlbum(count *int64, albumID uuid.UUID) error {
	args := s.Called(count, albumID)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) IsBandMemberAssociatedWithSong(songID uuid.UUID, bandMemberID uuid.UUID) (bool, error) {
	args := s.Called(songID, bandMemberID)
	return args.Bool(0), args.Error(1)
}

func (s *SongRepositoryMock) Create(song *model.Song) error {
	args := s.Called(song)
	return args.Error(0)
}

func (s *SongRepositoryMock) Update(song *model.Song) error {
	args := s.Called(song)
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

func (s *SongRepositoryMock) Delete(ids []uuid.UUID) error {
	args := s.Called(ids)
	return args.Error(0)
}

func (s *SongRepositoryMock) GetSettings(settings *model.SongSettings, settingsID uuid.UUID) error {
	args := s.Called(settings, settingsID)

	if len(args) > 1 {
		*settings = *args.Get(1).(*model.SongSettings)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) UpdateSettings(settings *model.SongSettings) error {
	args := s.Called(settings)
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

// Instruments

func (s *SongRepositoryMock) GetInstruments(instruments *[]model.Instrument, userID uuid.UUID) error {
	args := s.Called(instruments, userID)

	if len(args) > 1 {
		*instruments = *args.Get(1).(*[]model.Instrument)
	}

	return args.Error(0)
}

// Section Types

func (s *SongRepositoryMock) GetSectionTypes(sectionTypes *[]model.SongSectionType, userID uuid.UUID) error {
	args := s.Called(sectionTypes, userID)

	if len(args) > 1 {
		*sectionTypes = *args.Get(1).(*[]model.SongSectionType)
	}

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

func (s *SongRepositoryMock) DeleteSections(ids []uuid.UUID) error {
	args := s.Called(ids)
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
