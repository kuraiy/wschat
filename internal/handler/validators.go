package handler

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", validatePassword)
		v.RegisterValidation("username", validateUsername)
	}
}

func validatePassword(fl validator.FieldLevel) bool {
	var validChars = regexp.MustCompile(`^[A-Za-z0-9]{8,32}$`)
	var hasDigit = regexp.MustCompile(`[0-9]`)
	password := fl.Field().String()

	return validChars.MatchString(password) &&
		hasDigit.MatchString(password)
}

func validateUsername(fl validator.FieldLevel) bool {
	var validChars = regexp.MustCompile(`^[A-Za-z0-9]{8,32}$`)
	password := fl.Field().String()

	return validChars.MatchString(password)
}
