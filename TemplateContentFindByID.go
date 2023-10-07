package cms

func (cms *Cms) TemplateContentFindByID(templateID string) (string, error) {
	template, err := cms.TemplateFindByID(templateID)

	if err != nil {
		return "", err
	}

	if template == nil {
		return "", nil
	}

	content, err := template.GetString("content", "")

	if err != nil {
		return "", err
	}

	return content, nil
}
