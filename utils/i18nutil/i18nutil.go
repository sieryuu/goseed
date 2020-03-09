package i18nutil

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// GetTestingi18nBundle will return *i18n.Bundle for unit test
// Default language is English
func GetTestingi18nBundle() *i18n.Bundle {
	i18nBundle := i18n.NewBundle(language.English)
	i18nBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	i18nBundle.MustLoadMessageFile("../../../../utils/i18nutil/active.en.toml")
	return i18nBundle
}

// GetTestingLocalizer will return Localizer for unit test
func GetTestingLocalizer() *i18n.Localizer {
	return i18n.NewLocalizer(GetTestingi18nBundle(), "en")
}
