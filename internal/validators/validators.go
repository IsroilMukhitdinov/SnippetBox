package validators

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key string, value string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = value
	}
}

func (v *Validator) CheckField(ok bool, key string, value string) {
	if !ok {
		v.AddFieldError(key, value)
	}
}

func NotBlank(input string) bool {
	return strings.TrimSpace(input) != ""
}

func MaxLength(input string, length int) bool {
	return utf8.RuneCountInString(input) <= length
}

func PermittedValues(input int, values ...int) bool {
	for _, value := range values {
		if input == value {
			return true
		}
	}

	return false
}
