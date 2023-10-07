package cms

import (
	"log"

	"github.com/gouniverse/entitystore"
)

func (cms *Cms) TemplateFindByID(templateID string) (*entitystore.Entity, error) {
	return cms.EntityStore.EntityFindByID(templateID)
}

func (cms *Cms) TemplateContentFindByID(templateID string) string {
	template, err := cms.TemplateFindByID(templateID)

	if err != nil {
		log.Println("Template "+templateID+" error", err.Error())
		return ""
	}

	if template == nil {
		log.Println("Template " + templateID + " not found")
		return ""
	}

	content, err := template.GetString("content", "")

	if err != nil {
		log.Println("Template "+templateID+" content error", err.Error())
		return ""
	}

	return content
}
