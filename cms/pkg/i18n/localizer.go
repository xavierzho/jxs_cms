package i18n

import "golang.org/x/text/language"

type Localizer struct {
	Language language.Tag
	Bundle   *Bundle
}

func NewLocalizer(language language.Tag, bundle *Bundle) *Localizer {
	return &Localizer{
		Language: language,
		Bundle:   bundle,
	}
}

func (l *Localizer) LocalizeMessage(msg string) (string, bool) {
	message, ok := l.Bundle.messageTemplates[l.Language][msg]
	if ok {
		return message.Content, ok
	}
	return "", ok
}
