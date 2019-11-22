package validator

import (
	"reflect"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
)

var (
	validate = validator.New(&validator.Config{TagName: "validate"})
)

func Keyword(_ *validator.Validate, _ reflect.Value, _ reflect.Value, field reflect.Value, fieldType reflect.Type, _ reflect.Kind, _ string) bool {
	if fieldType.Kind() != reflect.String {
		return false
	}
	ok, _ := regexp.Match(`^\w+$`, []byte(field.Interface().(string)))
	return ok
}

func init() {
	err := validate.RegisterValidation("keyword", Keyword)
	if err != nil {
		panic(err)
	}
}

func RegisterGin() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("keyword", Keyword)
		if err != nil {
			panic(err)
		}
	} else {
		panic("Register gin validator's validation error")
	}
}

func Validate(value interface{}) error {
	err := validate.Struct(value)
	return err
}

func IsValidationError(err error) bool {
	switch err.(type) {
	case validator.ValidationErrors:
		return true
	}
	return false
}

func GenerateDetail(errors error) *map[string][]string {
	detail := map[string][]string{}
	for _, v := range errors.(validator.ValidationErrors) {
		if v.Param != "" {
			detail[v.Field] = append(detail[v.Field], v.Tag+"="+v.Param)
		} else {
			detail[v.Field] = append(detail[v.Field], v.Tag)
		}
	}
	return &detail
}
