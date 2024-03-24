package cms

func (cms *Cms) ShortcodesAdd(shortcodes []ShortcodeInterface) {
	for _, shortcode := range shortcodes {
		cms.ShortcodeAdd(shortcode)
	}
}
