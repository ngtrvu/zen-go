package zen

import (
	"errors"
)

const (
	DEFAULT_LOCALE string = "vi"
)

type AppError struct {
	Err          error
	Message      string
	Code         string
	Translations map[string]string
}

func NewAppError(code string, msg string) *AppError {
	var translations = make(map[string]string)
	return &AppError{
		Err:          errors.New(msg),
		Message:      msg,
		Code:         code,
		Translations: translations,
	}
}

func (e *AppError) AddTranslation(locale string, msg string) *AppError {
	e.Translations[locale] = msg

	return &AppError{
		Err:          e.Err,
		Message:      e.GetErrMsg(),
		Code:         e.Code,
		Translations: e.Translations,
	}
}

func (e *AppError) SetMessage(msg string) *AppError {
	e.Message = msg
	return &AppError{
		Err:          e.Err,
		Message:      e.GetErrMsg(),
		Code:         e.Code,
		Translations: e.Translations,
	}
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) GetErrMsg() string {
	val, ok := e.Translations[DEFAULT_LOCALE]
	if ok {
		return val
	}
	return e.Message
}

func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}
