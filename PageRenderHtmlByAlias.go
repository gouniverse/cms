package cms

import (
	"net/http"
	"strings"

	"github.com/gouniverse/cms/types"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/ui"
	"github.com/gouniverse/utils"
	"github.com/samber/lo"
)

// PageRenderHtmlByAlias builds the HTML of a page based on its alias
func (cms *Cms) PageRenderHtmlByAlias(r *http.Request, alias string, language string) string {
	page, err := cms.PageFindByAlias(alias)

	if err != nil {
		cms.logErrorWithContext(`At PageRenderHtmlByAlias`, err.Error())
		return hb.NewDiv().
			Text(`Page with alias '`).Text(alias).Text(`' not found`).
			ToHTML()
	}

	if page == nil {
		return hb.NewDiv().
			Text(`Page with alias '`).Text(alias).Text(`' not found`).
			ToHTML()
	}

	pageAttrs, err := page.GetAttributes()

	if err != nil {
		cms.logErrorWithContext(`At PageRenderHtmlByAlias`, err.Error())
		return hb.NewDiv().
			Text(`Page with alias '`).Text(alias).Text(`' io exception`).
			ToHTML()
	}

	pageContent := ""
	pageTitle := ""
	pageMetaKeywords := ""
	pageMetaDescription := ""
	pageMetaRobots := ""
	pageCanonicalURL := ""
	pageEditor := ""
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
		if attr.AttributeKey() == "editor" {
			pageEditor = attr.AttributeValue()
		}
	}

	if pageEditor == types.WEBPAGE_EDITOR_BLOCKEDITOR {
		if cms.blockEditorRenderer == nil {
			return "Block editor not configured"
		}

		if !utils.IsJSON(pageContent) {
			return "Malformed block content"
		}

		blocks, err := ui.UnmarshalJsonToBlocks(pageContent)

		if err != nil {
			return "Error parsing block content"
		}

		pageContent = cms.blockEditorRenderer(blocks)
	}

	if pageTemplateID == "" {
		return pageContent
	}

	finalContent := lo.If(pageTemplateID == "", pageContent).ElseF(func() string {
		content, err := cms.TemplateContentFindByID(pageTemplateID)
		if err != nil {
			cms.logErrorWithContext(`At PageRenderHtmlByAlias`, err.Error())
		}
		return content
	})

	html, err := cms.renderContentToHtml(r, finalContent, struct {
		PageContent         string
		PageCanonicalURL    string
		PageMetaDescription string
		PageMetaKeywords    string
		PageMetaRobots      string
		PageTitle           string
		Language            string
	}{
		PageContent:         pageContent,
		PageCanonicalURL:    pageCanonicalURL,
		PageMetaDescription: pageMetaDescription,
		PageMetaKeywords:    pageMetaKeywords,
		PageMetaRobots:      pageMetaRobots,
		PageTitle:           pageTitle,
	})

	if err != nil {
		cms.logErrorWithContext(`At PageRenderHtmlByAlias`, err.Error())
		return hb.NewDiv().Text(`error occurred`).ToHTML()
	}

	return html
}

// renderContentToHtml renders the content to HTML
//
// This is done in the following steps (sequence is important):
// 1. replaces placeholders with values
// 2. renders the blocks
// 3. renders the shortcodes
// 3. renders the translations
// 4. returns the HTML
//
// Parameters:
// - r: the HTTP request
// - content: the content to render
// - options: the options for the rendering
//
// Returns:
// - html: the rendered HTML
// - err: the error, if any, or nil otherwise
func (cms *Cms) renderContentToHtml(r *http.Request, content string, options struct {
	PageContent         string
	PageCanonicalURL    string
	PageMetaDescription string
	PageMetaKeywords    string
	PageMetaRobots      string
	PageTitle           string
	Language            string
}) (html string, err error) {
	replacements := map[string]string{
		"PageContent":         options.PageContent,
		"PageCanonicalUrl":    options.PageCanonicalURL,
		"PageMetaDescription": options.PageMetaDescription,
		"PageMetaKeywords":    options.PageMetaKeywords,
		"PageRobots":          options.PageMetaRobots,
		"PageTitle":           options.PageTitle,
	}

	for key, value := range replacements {
		content = strings.ReplaceAll(content, "[["+key+"]]", value)
		content = strings.ReplaceAll(content, "[[ "+key+" ]]", value)
	}

	content, err = cms.ContentRenderBlocks(content)

	if err != nil {
		return "", err
	}

	content, err = cms.ContentRenderShortcodes(r, content)

	if err != nil {
		return "", err
	}

	language := lo.If(options.Language == "", "en").Else(options.Language)

	content, err = cms.ContentRenderTranslations(content, language)

	if err != nil {
		return "", err
	}

	return content, nil
}
