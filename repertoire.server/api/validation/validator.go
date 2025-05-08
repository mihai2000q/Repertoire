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

func (v *Validator) Validate(request any) *wrapper.ErrorCode {
	err := v.validate.Struct(request)
	if err != nil {
		return wrapper.BadRequestError(err)
	}
	return nil
}

func registerCustomValidators(validate *validator.Validate) error {
	err := validate.RegisterValidation("has_upper", HasUpper)
	if err != nil {
		return err
	}

	err = validate.RegisterValidation("has_lower", HasLower)
	if err != nil {
		return err
	}

	err = validate.RegisterValidation("has_digit", HasDigit)
	if err != nil {
		return err
	}

	err = validate.RegisterValidation("difficulty_enum", DifficultyEnum)
	if err != nil {
		return err
	}

	err = validate.RegisterValidation("search_type_enum", SearchTypeEnum)
	if err != nil {
		return err
	}

	err = validate.RegisterValidation("youtube_link", YoutubeLink)
	if err != nil {
		return err
	}

	err = validate.RegisterValidation("color", Color)
	if err != nil {
		return err
	}

	err = validate.RegisterValidation("order_by", OrderBy)
	if err != nil {
		return err
	}

	err = validate.RegisterValidation("search_by", SearchBy)
	if err != nil {
		return err
	}

	err = validate.RegisterValidation("notblank", validators.NotBlank)
	if err != nil {
		return err
	}
	return nil
}
