package cms

import "strings"

// ContentRenderTranslations renders the translations in a string
func (cms *Cms) ContentRenderTranslations(content string, language string) (string, error) {
	translationIDs := ContentFindIdsByPatternPrefix(content, "TRANSLATION")

	var err error
	for _, translationID := range translationIDs {
		content, err = cms.ContentRenderTranslationByIdOrHandle(content, translationID, language)

		if err != nil {
			return content, err
		}
	}

	return content, nil
}

// ContentRenderTranslationByIdOrHandle renders the translation specified by the ID in a content
// if the blockID is empty or not found the initial content is returned
func (cms *Cms) ContentRenderTranslationByIdOrHandle(content string, translationID string, language string) (string, error) {
	if translationID == "" {
		return content, nil
	}

	translation, err := cms.TranslationFindByIdOrHandle(translationID, language)

	if err != nil {
		return "", err
	}

	content = strings.ReplaceAll(content, "[[TRANSLATION_"+translationID+"]]", translation)
	content = strings.ReplaceAll(content, "[[ TRANSLATION_"+translationID+" ]]", translation)

	return content, nil
}
