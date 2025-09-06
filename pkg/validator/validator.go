package validator

import (
    "github.com/go-playground/validator/v10"
)

type Validator interface {
    Validate(i interface{}) error
}

type validatorImpl struct {
    validator *validator.Validate
}

func NewValidator() Validator {
    return &validatorImpl{
        validator: validator.New(),
    }
}

func (v *validatorImpl) Validate(i interface{}) error {
    return v.validator.Struct(i)
}