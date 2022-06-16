package domain

import validation "github.com/go-ozzo/ozzo-validation/v4"

var CV = &CustomValidator{}

// Id...chi.URLParamでparameterをGetするとき、stringになり、型を一定のものにして副作用がないようにするためにこれを利用する
type Id string

type CustomValidator struct{}

func (cv *CustomValidator) Validate(i interface{}) error {
	if c, ok := i.(validation.Validatable); ok {
		return c.Validate()
	}
	return nil
}
