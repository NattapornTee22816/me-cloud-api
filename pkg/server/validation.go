package server

import "github.com/go-playground/validator/v10"

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Engine() any {
	//TODO implement me
	panic("implement me")
}

func (v *structValidator) ValidateStruct(out any) error {
	return v.validate.Struct(out)
}
