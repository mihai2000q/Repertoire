package service

import (
	"mime/multipart"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/user"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type UserService interface {
	DeleteProfilePicture(token string) *wrapper.ErrorCode
	Get(id uuid.UUID) (user model.User, e *wrapper.ErrorCode)
	SaveProfilePicture(file *multipart.FileHeader, token string) *wrapper.ErrorCode
	Update(request requests.UpdateUserRequest, token string) *wrapper.ErrorCode
}

type userService struct {
	deleteProfilePictureFromUser user.DeleteProfilePictureFromUser
	getUser                      user.GetUser
	saveProfilePictureToUser     user.SaveProfilePictureToUser
	updateUser                   user.UpdateUser
}

func NewUserService(
	deleteProfilePictureFromUser user.DeleteProfilePictureFromUser,
	getUser user.GetUser,
	saveProfilePictureToUser user.SaveProfilePictureToUser,
	updateUser user.UpdateUser,
) UserService {
	return &userService{
		deleteProfilePictureFromUser: deleteProfilePictureFromUser,
		getUser:                      getUser,
		saveProfilePictureToUser:     saveProfilePictureToUser,
		updateUser:                   updateUser,
	}
}

func (u *userService) DeleteProfilePicture(token string) *wrapper.ErrorCode {
	return u.deleteProfilePictureFromUser.Handle(token)
}

func (u *userService) Get(id uuid.UUID) (model.User, *wrapper.ErrorCode) {
	return u.getUser.Handle(id)
}

func (u *userService) SaveProfilePicture(file *multipart.FileHeader, token string) *wrapper.ErrorCode {
	return u.saveProfilePictureToUser.Handle(file, token)
}

func (u *userService) Update(request requests.UpdateUserRequest, token string) *wrapper.ErrorCode {
	return u.updateUser.Handle(request, token)
}
