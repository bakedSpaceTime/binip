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
	body := m.form.View()

	// Show error or status message if present
	if m.msg != "" {
		body = body + "\n\n" + styles.ErrorStyle.Render(m.msg)
	}

	return body
}

func (m *mainModel) operationalView() string {
	switch m.operationalMode {
	case listView:
		return m.listView()
	case detailView:
		return m.detailView()
	case createView, editView, deleteConfirmView:
		// Form-based views
		body := m.form.View()
		if m.msg != "" {
			body = body + "\n\n" + styles.StatusStyle.Render(m.msg)
		}
		return body
	default:
		// Fallback
		return m.db.String()
	}
}

// listView shows the list of records
func (m *mainModel) listView() string {
	// TODO: Implement when CRUD UI is designed
	// For now, show the database info as before
	body := m.db.String()

	if m.msg != "" {
		body = body + "\n\n" + styles.StatusStyle.Render(m.msg)
	}

	return body
}

// detailView shows details of a single record
func (m *mainModel) detailView() string {
	// TODO: Implement when CRUD UI is designed
	body := "Detail view for record: " + m.currentRecordID

	if m.msg != "" {
		body = body + "\n\n" + styles.StatusStyle.Render(m.msg)
	}

	return body
}
