package validation

import (
	"repertoire/utils"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() Validator {
	validate := validator.New(validator.WithRequiredStructEnabled())
	registerCustomValidators(validate)
	return Validator{
		validate: validate,
	}
}

func (v *Validator) Validate(request interface{}) *utils.ErrorCode {
	err := v.validate.Struct(request)
	if err != nil {
		return utils.BadRequestError(err)
	}
	return nil
}

func registerCustomValidators(validate *validator.Validate) {
	validate.RegisterValidation("notblank", validators.NotBlank)
}
