package auth

import (
	"errors"
	"github.com/google/uuid"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"strings"
)

type SignUp struct {
	jwtService     service.JwtService
	bCryptService  service.BCryptService
	userRepository repository.UserRepository
}

func NewSignUp(
	jwtService service.JwtService,
	bCryptService service.BCryptService,
	userRepository repository.UserRepository,
) SignUp {
	return SignUp{
		jwtService:     jwtService,
		bCryptService:  bCryptService,
		userRepository: userRepository,
	}
}

func (s *SignUp) Handle(request requests.SignUpRequest) (string, *wrapper.ErrorCode) {
	var user model.User

	// check if the user already exists
	email := strings.ToLower(request.Email)
	err := s.userRepository.GetByEmail(&user, email)
	if err != nil {
		return "", wrapper.InternalServerError(err)
	}
	if !reflect.ValueOf(user).IsZero() {
		return "", wrapper.BadRequestError(errors.New("user already exists"))
	}

	// hash the password
	hashedPassword, err := s.bCryptService.Hash(request.Password)
	if err != nil {
		return "", wrapper.InternalServerError(err)
	}

	// create user
	user = model.User{
		ID:       uuid.New(),
		Name:     request.Name,
		Email:    email,
		Password: hashedPassword,
	}
	s.createAndAttachDefaultData(&user)
	err = s.userRepository.Create(&user)
	if err != nil {
		return "", wrapper.InternalServerError(err)
	}

	return s.jwtService.CreateToken(user)
}

func (s *SignUp) createAndAttachDefaultData(user *model.User) {
	var guitarTunings []model.GuitarTuning
	var songSectionTypes []model.SongSectionType
	var bandMemberRoles []model.BandMemberRole

	for i, guitarTuning := range model.DefaultGuitarTunings {
		guitarTunings = append(guitarTunings, model.GuitarTuning{
			ID:     uuid.New(),
			Name:   guitarTuning,
			Order:  uint(i),
			UserID: user.ID,
		})
	}

	for i, songSectionType := range model.DefaultSongSectionTypes {
		songSectionTypes = append(songSectionTypes, model.SongSectionType{
			ID:     uuid.New(),
			Name:   songSectionType,
			Order:  uint(i),
			UserID: user.ID,
		})
	}

	for i, songSectionType := range model.DefaultSongSectionTypes {
		songSectionTypes = append(songSectionTypes, model.SongSectionType{
			ID:     uuid.New(),
			Name:   songSectionType,
			Order:  uint(i),
			UserID: user.ID,
		})
	}

	for i, bandMemberRole := range model.DefaultBandMemberRoles {
		bandMemberRoles = append(bandMemberRoles, model.BandMemberRole{
			ID:     uuid.New(),
			Name:   bandMemberRole,
			Order:  uint(i),
			UserID: user.ID,
		})
	}

	user.GuitarTunings = guitarTunings
	user.SongSectionTypes = songSectionTypes
	user.BandMemberRoles = bandMemberRoles
}
