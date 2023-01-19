package validator

import (
	validators "github.com/go-playground/validator/v10"
)

type Validator interface {
	ValidateStruct(inf interface{}) error
}

type validator struct {
	validator *validators.Validate
}

func New() Validator {
	v := validators.New()
	return &validator{
		validator: v,
	}
}

func (v *validator) ValidateStruct(inf interface{}) error {

	return v.validator.Struct(inf)
}
