package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Todo struct {
	Model
	Title       string `gorm:"title" json:"title"`
	Description string `gorm:"description" json:"description"`
	Completed   bool   `gorm:"completed" json:"completed"`
}

func (a Todo) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(
			&a.Title,
			validation.Required.Error("is required"),
			validation.RuneLength(1, 50).Error("size is 1～50"),
		),
		validation.Field(
			&a.Description,
			validation.Required.Error("is required"),
			validation.RuneLength(1, 100).Error("size is 1～100"),
		),
	)
}
