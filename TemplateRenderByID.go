package cms

import (
	"errors"
	"net/http"
)

// TemplateRenderHtmlByID builds the HTML of a template based on its ID
func (cms *Cms) TemplateRenderHtmlByID(r *http.Request, templateID string, options struct {
	PageContent         string
	PageCanonicalURL    string
	PageMetaDescription string
	PageMetaKeywords    string
	PageMetaRobots      string
	PageTitle           string
	Language            string
}) (string, error) {
	if templateID == "" {
		return "", errors.New("template id is empty")
	}

	content, err := cms.TemplateContentFindByID(templateID)

	if err != nil {
		return "", err
	}

	html, err := cms.renderContentToHtml(r, content, options)

	if err != nil {
		return "", err
	}

	return html, nil
}
