package validator

import (
	"github.com/go-playground/validator/v10"
)

var v = validator.New()

func ValidateStruct(s any) error { return v.Struct(s) }
