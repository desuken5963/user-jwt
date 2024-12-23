package utils

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct based on `validate` tags.
func ValidateStruct(s interface{}) map[string]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		// JSONタグを取得
		field, _ := reflect.TypeOf(s).Elem().FieldByName(err.Field())
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = strings.ToLower(err.Field())
		}

		errors[jsonTag] = err.Tag() // バリデーションルールを格納
	}
	return errors
}
