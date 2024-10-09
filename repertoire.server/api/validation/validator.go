package validation

import (
	"github.com/go-playground/validator/v10"
	"repertoire/utils"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() Validator {
	return Validator{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (v *Validator) Validate(request interface{}) *utils.ErrorCode {
	err := v.validate.Struct(request)
	if err != nil {
		return utils.BadRequestError(err)
	}
	return nil
}
