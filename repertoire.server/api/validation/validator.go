package validation

import (
	"context"
	"repertoire/server/internal/httperror"

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
		if err := registerCustomValidators(validate); err != nil {
			panic(err)
		}
	}

	return &Validator{
		validate: validate,
	}
}

func (v *Validator) Validate(request any) *httperror.ErrorCode {
	if err := v.validate.Struct(request); err != nil {
		return httperror.BadRequestError(err)
	}
	return nil
}

func registerCustomValidators(validate *validator.Validate) error {
	if err := validate.RegisterValidation("has_upper", HasUpper); err != nil {
		return err
	}

	if err := validate.RegisterValidation("has_lower", HasLower); err != nil {
		return err
	}

	if err := validate.RegisterValidation("has_digit", HasDigit); err != nil {
		return err
	}

	if err := validate.RegisterValidation("difficulty_enum", DifficultyEnum); err != nil {
		return err
	}

	if err := validate.RegisterValidation("search_type_enum", SearchTypeEnum); err != nil {
		return err
	}

	if err := validate.RegisterValidation("youtube_link", YoutubeLink); err != nil {
		return err
	}

	if err := validate.RegisterValidation("color", Color); err != nil {
		return err
	}

	if err := validate.RegisterValidation("order_by", OrderBy); err != nil {
		return err
	}

	if err := validate.RegisterValidation("search_order", SearchOrder); err != nil {
		return err
	}

	if err := validate.RegisterValidation("search_by", SearchBy); err != nil {
		return err
	}

	if err := validate.RegisterValidation("search_filter", SearchFilter); err != nil {
		return err
	}

	if err := validate.RegisterValidation("unique_ids", UniqueIDs); err != nil {
		return err
	}

	if err := validate.RegisterValidation("notblank", validators.NotBlank); err != nil {
		return err
	}
	return nil
}
