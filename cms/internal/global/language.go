package global

import (
	"golang.org/x/text/language"
)

var Language = language.Chinese

func SetLanguage(languageStr string) (err error) {
	language, err := language.Parse(languageStr)
	if err != nil {
		return err
	}
	Language = language

	return nil
}
