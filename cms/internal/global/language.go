package global

import (
	"golang.org/x/text/language"
)

var Language = language.Chinese

func SetLanguage(languageStr string) (err error) {
	Language, err = language.Parse(languageStr)
	if err != nil {
		return err
	}

	return nil
}
