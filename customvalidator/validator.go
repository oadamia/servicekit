package customvalidator

import "github.com/go-playground/validator/v10"

type Custom struct {
	validator *validator.Validate
}

func (cv *Custom) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func New() *Custom {
	return &Custom{validator: validator.New()}
}
