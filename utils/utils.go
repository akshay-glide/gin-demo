package utils

import (
	"errors"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type ErrStructure struct {
	Field string `json:"field"`
	Param string `json:"param"`
	Tag   string `json:"err"`
	Info  string `json:"info"`
}

func ValidateStruct(structVar any) (any, error) {
	validate := validator.New()
	if err := validate.RegisterValidation("regex", validateRegex); err != nil {
		log.Err(err).Msg("register failed")
		return "", err
	}

	err := validate.Struct(structVar)
	if err != nil {
		msg := "Validation failed for field(s): "
		fields := []string{}
		errStructs := []ErrStructure{}

		for _, err := range err.(validator.ValidationErrors) {

			errStruct, regexErrMessgae := getRegexErrorMessage(err.Field(), err.Param(), err.Tag())

			fields = append(fields, strings.ToLower(err.Field())+regexErrMessgae)

			errStructs = append(errStructs, errStruct)
		}
		msg = msg + strings.Join(fields, "; ")

		return errStructs, errors.New(msg)
	}

	return "", nil

}

func validateRegex(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	regex := fl.Param()
	match, _ := regexp.MatchString(regex, value)
	return match
}

func getRegexErrorMessage(field, param, tag string) (ErrStructure, string) {

	var errMessage string

	errMessage = `field ` + field + ` doesn't present or contains unwanted charachters`

	switch tag {
	case "max":
		errMessage = errMessage + ` field ` + field + ` can't have more than ` + param + ` charachters`
	case "min":
		errMessage = errMessage + ` field ` + field + ` should have mminimum ` + param + ` charachters`
	}

	return ErrStructure{
		Field: field,
		Param: param,
		Tag:   tag,
		Info:  errMessage,
	}, errMessage
}

func GetRandomAlphaNumeric(length int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	sb := strings.Builder{}
	sb.Grow(length)
	for i := 0; i < length; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	res := sb.String()
	return res
}
