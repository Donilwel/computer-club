package utils

import (
	"errors"
)

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		return errors.New("invalid input data")
	}
	return nil
}
