package validation

import (
	"context"
	"repertoire/server/internal/wrapper"

	"go.uber.org/fx"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator(lc fx.Lifecycle) *Validator {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if lc != nil { // Null on Unit Testing
		lc.Append(fx.Hook{
			OnStart: func(context.Context) error {
				return registerCustomValidators(validate)
			},
		})
	} else {
		err := registerCustomValidators(validate)
		if err != nil {
			panic(err)
		}
	}

	return &Validator{
		validate: validate,
	}
}

func (v *Validator) Validate(request interface{}) *wrapper.ErrorCode {
	err := v.validate.Struct(request)
	if err != nil {
		return wrapper.BadRequestError(err)
	}
	return nil
}

func registerCustomValidators(validate *validator.Validate) error {
	err := validate.RegisterValidation("hasUpper", HasUpper)
	if err != nil {
		return err
	}

	err = validate.RegisterValidation("hasLower", HasLower)
	if err != nil {
		return err
	}

	err = validate.RegisterValidation("hasDigit", HasDigit)
	if err != nil {
		return err
	}

	err = validate.RegisterValidation("isDifficultyEnum", IsDifficultyEnum)
	if err != nil {
		return err
	}

	err = validate.RegisterValidation("isYoutubeLink", IsYoutubeLink)
	if err != nil {
		return err
	}

	err = validate.RegisterValidation("isColor", IsColor)
	if err != nil {
		return err
	}

	err = validate.RegisterValidation("notblank", validators.NotBlank)
	if err != nil {
		return err
	}
	return nil
}
