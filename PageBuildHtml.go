package cms

import (
	"net/http"
	"strings"

	"github.com/samber/lo"
)

// PageRenderHtmlByAlias builds the HTML of a page based on its alias
func (cms *Cms) PageRenderHtmlByAlias(r *http.Request, alias string, language string) string {
	page, err := cms.PageFindByAlias(alias)

	if err != nil {
		cms.LogStore.ErrorWithContext("At pageBuldHTMLByAlias", err.Error())
		return "Page with alias '" + alias + "' not found"
	}

	if page == nil {
		return "Page with alias '" + alias + "' not found"
	}

	pageAttrs, err := page.GetAttributes()

	if err != nil {
		return "Page '" + alias + "' io exception"
	}

	pageContent := ""
	pageTitle := ""
	pageMetaKeywords := ""
	pageMetaDescription := ""
	pageMetaRobots := ""
	pageCanonicalURL := ""
	pageTemplateID := ""
	for _, attr := range pageAttrs {
		if attr.AttributeKey() == "content" {
			pageContent = attr.AttributeValue()
		}
		if attr.AttributeKey() == "title" {
			pageTitle = attr.AttributeValue()
		}
		if attr.AttributeKey() == "meta_keywords" {
			pageMetaKeywords = attr.AttributeValue()
		}
		if attr.AttributeKey() == "meta_description" {
			pageMetaDescription = attr.AttributeValue()
		}
		if attr.AttributeKey() == "meta_robots" {
			pageMetaRobots = attr.AttributeValue()
		}
		if attr.AttributeKey() == "canonical_url" {
			pageCanonicalURL = attr.AttributeValue()
		}
		if attr.AttributeKey() == "template_id" {
			pageTemplateID = attr.AttributeValue()
		}
	}

	finalContent := lo.If(pageTemplateID == "", pageContent).ElseF(func() string {
		content, err := cms.TemplateContentFindByID(pageTemplateID)
		if err != nil {
			cms.LogStore.ErrorWithContext("At pageBuldHTMLByAlias", err.Error())
		}
		return content
	})

	replacements := map[string]string{
		"PageContent":         pageContent,
		"PageCanonicalUrl":    pageCanonicalURL,
		"PageMetaDescription": pageMetaDescription,
		"PageMetaKeywords":    pageMetaKeywords,
		"PageRobots":          pageMetaRobots,
		"PageTitle":           pageTitle,
	}

	for key, value := range replacements {
		finalContent = strings.ReplaceAll(finalContent, "[["+key+"]]", value)
		finalContent = strings.ReplaceAll(finalContent, "[[ "+key+" ]]", value)
	}

	finalContent, err = cms.ContentRenderBlocks(finalContent)

	if err != nil {
		cms.LogStore.ErrorWithContext("At pageBuldHTMLByAlias", err.Error())
	}

	finalContent, err = cms.ContentRenderShortcodes(r, finalContent)

	if err != nil {
		cms.LogStore.ErrorWithContext("At pageBuldHTMLByAlias", err.Error())
	}

	finalContent, err = cms.ContentRenderTranslations(finalContent, language)

	if err != nil {
		cms.LogStore.ErrorWithContext("At pageBuldHTMLByAlias", err.Error())
	}

	return finalContent
}
