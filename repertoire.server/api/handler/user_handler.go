package handler

import (
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/api/server"
	"repertoire/server/api/validation"
	"repertoire/server/domain/provider"
	"repertoire/server/domain/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service             service.UserService
	currentUserProvider provider.CurrentUserProvider
	server.BaseHandler
}

func NewUserHandler(
	service service.UserService,
	currentUserProvider provider.CurrentUserProvider,
	validator *validation.Validator,
) *UserHandler {
	return &UserHandler{
		service:             service,
		currentUserProvider: currentUserProvider,
		BaseHandler: server.BaseHandler{
			Validator: validator,
		},
	}
}

func (u UserHandler) GetCurrentUser(c *gin.Context) {
	token := u.GetTokenFromContext(c)

	user, errorCode := u.currentUserProvider.Get(token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u UserHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, errorCode := u.service.Get(id)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u UserHandler) SignUp(c *gin.Context) {
	var request requests.SignUpRequest
	errCode := u.BindAndValidate(c, &request)
	if errCode != nil {
		_ = c.AbortWithError(errCode.Code, errCode.Error)
		return
	}

	token, errCode := u.service.SignUp(request)
	if errCode != nil {
		_ = c.AbortWithError(errCode.Code, errCode.Error)
		return
	}

	c.JSON(http.StatusOK, token)
}

func (u UserHandler) Update(c *gin.Context) {
	var request requests.UpdateUserRequest
	errCode := u.BindAndValidate(c, &request)
	if errCode != nil {
		_ = c.AbortWithError(errCode.Code, errCode.Error)
		return
	}

	token := u.GetTokenFromContext(c)

	errCode = u.service.Update(request, token)
	if errCode != nil {
		_ = c.AbortWithError(errCode.Code, errCode.Error)
		return
	}

	u.SendMessage(c, "user has been updated successfully!")
}

func (u UserHandler) Delete(c *gin.Context) {
	token := u.GetTokenFromContext(c)

	errCode := u.service.Delete(token)
	if errCode != nil {
		_ = c.AbortWithError(errCode.Code, errCode.Error)
		return
	}

	u.SendMessage(c, "user has been deleted successfully!")
}

// Pictures

func (u UserHandler) SaveProfilePicture(c *gin.Context) {
	file, err := c.FormFile("profile_pic")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	token := u.GetTokenFromContext(c)

	errorCode := u.service.SaveProfilePicture(file, token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	u.SendMessage(c, "profile picture has been saved to current user successfully!")
}

func (u UserHandler) DeleteProfilePicture(c *gin.Context) {
	token := u.GetTokenFromContext(c)

	errorCode := u.service.DeleteProfilePicture(token)
	if errorCode != nil {
		_ = c.AbortWithError(errorCode.Code, errorCode.Error)
		return
	}

	u.SendMessage(c, "profile picture has been deleted from current user successfully")
}
