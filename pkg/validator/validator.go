// pkg/validator/validator.go

package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidators(v *validator.Validate) {
    // 手机号验证
    v.RegisterValidation("phone", validatePhone)
}

func validatePhone(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    reg := regexp.MustCompile(`^1[3-9]\d{9}$`)
    return reg.MatchString(phone)
}