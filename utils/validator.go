package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	var messages []string
	for _, fe := range validationErrors {
		field := fe.Field()
		switch fe.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("%s wajib diisi", field))
		case "email":
			messages = append(messages, fmt.Sprintf("%s harus format email yang valid", field))
		case "min":
			messages = append(messages, fmt.Sprintf("%s minimal %s karakter", field, fe.Param()))
		case "max":
			messages = append(messages, fmt.Sprintf("%s maksimal %s karakter", field, fe.Param()))
		case "gt":
			messages = append(messages, fmt.Sprintf("%s harus lebih dari %s", field, fe.Param()))
		case "gte":
			messages = append(messages, fmt.Sprintf("%s minimal %s", field, fe.Param()))
		case "oneof":
			messages = append(messages, fmt.Sprintf("%s harus salah satu dari: %s", field, fe.Param()))
		default:
			messages = append(messages, fmt.Sprintf("%s tidak valid", field))
		}
	}

	return fmt.Errorf("%s", strings.Join(messages, "; "))
}
