package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func Init() error {
	bundle = i18n.NewBundle(language.Japanese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	localesDir := "locales"
	files, err := os.ReadDir(localesDir)
	if err != nil {
		return fmt.Errorf("failed to read locales directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			path := filepath.Join(localesDir, file.Name())
			bundle.MustLoadMessageFile(path)
		}
	}

	return nil
}

func GetMessage(lang, messageID string, templateData map[string]interface{}) string {
	localizer := i18n.NewLocalizer(bundle, lang)
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
	if err != nil {
		// Fallback to default language (Japanese) if translation is not found
		localizer = i18n.NewLocalizer(bundle, language.Japanese.String())
		msg, _ = localizer.Localize(&i18n.LocalizeConfig{
			MessageID:    messageID,
			TemplateData: templateData,
		})
	}
	return msg
}
