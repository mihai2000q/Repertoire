package auth

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/model"
	"repertoire/utils/wrapper"
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
	if user.ID == uuid.Nil {
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
	err = s.userRepository.Create(&user)
	s.createAndAttachDefaultData(&user)
	if err != nil {
		return "", wrapper.InternalServerError(err)
	}

	return s.jwtService.CreateToken(user)
}

var defaultGuitarTuning = []string{
	"E Standard", "Eb Standard", "D Standard", "C# Standard", "C Standard", "B Standard", "A# Standard", "A Standard",
	"Drop D", "Drop C#", "Drop C", "Drop B", "Drop A#", "Drop A",
}
var defaultSongSectionTypes = []string{"Intro", "Verse", "Chorus", "Interlude", "Breakdown", "Solo", "Riff", "Outro"}

func (s *SignUp) createAndAttachDefaultData(user *model.User) {
	var guitarTunings []model.GuitarTuning
	var songSectionTypes []model.SongSectionType

	for _, guitarTuning := range defaultGuitarTuning {
		guitarTunings = append(guitarTunings, model.GuitarTuning{
			ID:     uuid.New(),
			Name:   guitarTuning,
			UserID: user.ID,
		})
	}

	for _, songSectionType := range defaultSongSectionTypes {
		songSectionTypes = append(songSectionTypes, model.SongSectionType{
			ID:     uuid.New(),
			Name:   songSectionType,
			UserID: user.ID,
		})
	}

	user.GuitarTunings = guitarTunings
	user.SongSectionTypes = songSectionTypes
}
