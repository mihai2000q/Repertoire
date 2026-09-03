package service

import (
	"mime/multipart"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/user"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type UserService interface {
	Delete(token string) *httperror.ErrorCode
	DeleteProfilePicture(token string) *httperror.ErrorCode
	Get(id uuid.UUID) (user model.User, e *httperror.ErrorCode)
	SaveProfilePicture(file *multipart.FileHeader, token string) *httperror.ErrorCode
	SignUp(request requests.SignUpRequest) (string, *httperror.ErrorCode)
	Update(request requests.UpdateUserRequest, token string) *httperror.ErrorCode
}

type userService struct {
	deleteUser                   user.DeleteUser
	deleteProfilePictureFromUser user.DeleteProfilePictureFromUser
	getUser                      user.GetUser
	saveProfilePictureToUser     user.SaveProfilePictureToUser
	signUp                       user.SignUp
	updateUser                   user.UpdateUser
}

func NewUserService(
	deleteUser user.DeleteUser,
	deleteProfilePictureFromUser user.DeleteProfilePictureFromUser,
	getUser user.GetUser,
	saveProfilePictureToUser user.SaveProfilePictureToUser,
	signUp user.SignUp,
	updateUser user.UpdateUser,
) UserService {
	return &userService{
		deleteUser:                   deleteUser,
		deleteProfilePictureFromUser: deleteProfilePictureFromUser,
		getUser:                      getUser,
		saveProfilePictureToUser:     saveProfilePictureToUser,
		signUp:                       signUp,
		updateUser:                   updateUser,
	}
}

func (u *userService) Delete(token string) *httperror.ErrorCode {
	return u.deleteUser.Handle(token)
}

func (u *userService) DeleteProfilePicture(token string) *httperror.ErrorCode {
	return u.deleteProfilePictureFromUser.Handle(token)
}

func (u *userService) Get(id uuid.UUID) (model.User, *httperror.ErrorCode) {
	return u.getUser.Handle(id)
}

func (u *userService) SaveProfilePicture(file *multipart.FileHeader, token string) *httperror.ErrorCode {
	return u.saveProfilePictureToUser.Handle(file, token)
}

func (u *userService) SignUp(request requests.SignUpRequest) (string, *httperror.ErrorCode) {
	return u.signUp.Handle(request)
}

func (u *userService) Update(request requests.UpdateUserRequest, token string) *httperror.ErrorCode {
	return u.updateUser.Handle(request, token)
}
