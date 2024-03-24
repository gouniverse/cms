package cms

func (cms *Cms) ShortcodeAdd(shortcode ShortcodeInterface) {
	cms.shortcodes = append(cms.shortcodes, shortcode)
}
