package dto

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"time"
)

type SearchDto struct {
	Location    string    `json:"location"`
	GuestNumber int       `json:"guest_number" validate:"min=2"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end" validate:"gtefield=Start"`
}

func ValidateStruct(dto SearchDto) error {
	var validate = validator.New()
	err := validate.Struct(dto)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return fmt.Errorf("invalid validation error: %w", err)
		}

		for _, err := range err.(validator.ValidationErrors) {
			// Print each validation error
			fmt.Printf("Field '%s' failed validation for tag '%s'\n", err.StructNamespace(), err.Tag())
		}
	}

	return err
}
