package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func Validate(data any) error {
	return validate.Struct(data)
}

// Key reports the first violated rule as "Field.tag", letting callers map a
// rule onto their own error without importing the validator package themselves.
// When several rules fail, the earliest struct field wins.
func Key(err error) (string, bool) {
	var fieldErrors validator.ValidationErrors
	if !errors.As(err, &fieldErrors) || len(fieldErrors) == 0 {
		return "", false
	}

	return fieldErrors[0].Field() + "." + fieldErrors[0].Tag(), true
}
