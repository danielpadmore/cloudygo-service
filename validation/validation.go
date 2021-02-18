package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/danielpadmore/cloudygo-service/logs"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

// Validator contains interfaces to validate and translate validation messages
type Validator struct {
	Validate *validator.Validate
	trans    ut.Translator
}

// New registers custom validations and creates a validator
func New(logger logs.Logger) Validator {
	translator := en.New()
	uni := ut.New(translator, translator)

	trans, found := uni.GetTranslator("en")
	if !found {
		logger.Fatal(logs.NewLog("SETUP", fmt.Sprintf("EN translator unfound")))
	}

	val := validator.New()

	if err := en_translations.RegisterDefaultTranslations(val, trans); err != nil {
		logger.Fatal(logs.NewLog("SETUP", fmt.Sprintf("Error registering default EN translation: %s", err.Error())))
	}

	_ = val.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	val.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return Validator{val, trans}
}

// ConcatReasons returns a string with all validation reasons separated by commas
func (val *Validator) ConcatReasons(err error) string {
	valErrs := err.(validator.ValidationErrors)
	reasons := make([]string, 0, len(valErrs))
	for _, e := range valErrs {
		reasons = append(reasons, e.Translate(val.trans))
	}
	return strings.Join(reasons, ", ")
}
