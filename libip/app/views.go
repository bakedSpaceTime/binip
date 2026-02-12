package app

import (
	"github.com/bakedSpaceTime/binip/libip/styles"
)

func (m *mainModel) footerView() string {
	helpView := m.help.View(m.keys)
	// Add debug status to expanded help view
	if m.help.ShowAll {
		debugStatus := "off"
		if m.config.Debug {
			debugStatus = "on"
		}
		helpView += "\n" + m.help.Styles.FullDesc.Render("Debug: ") + m.help.Styles.FullKey.Render(debugStatus)
	}
	return styles.FooterStyle.
		Render(helpView)
}

func (m *mainModel) onboardingView() string {
	return m.form.View()
}
