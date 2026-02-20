package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
)

// handleMessage routes messages to the appropriate state handler
func (m *mainModel) handleMessage(msg tea.Msg) tea.Cmd {
	switch m.state {
	case onboarding:
		return m.handleOnboardingMessage(msg)
	case operational:
		return m.handleOperationalMessage(msg)
	}
	return nil
}

// handleOnboardingMessage handles messages during onboarding state
func (m *mainModel) handleOnboardingMessage(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case prefixSelectedMsg:
		if m.config.Debug {
			spew.Fdump(m.config.DebugWriter, msg, "prefix selected")
		}

		if msg.isCustom {
			// User wants custom prefix, go to custom entry form
			return m.transitionToOnboardingState(enteringCustomPrefix, "")
		}
		// User selected preset, go to confirmation
		return m.transitionToOnboardingState(confirmingPrefix, msg.prefix)

	case customPrefixEnteredMsg:
		if m.config.Debug {
			spew.Fdump(m.config.DebugWriter, msg, "custom prefix entered")
		}

		if msg.valid {
			// Valid custom prefix entered, go to confirmation
			return m.transitionToOnboardingState(confirmingPrefix, msg.prefix)
		}
		// Invalid prefix, validation failed - stay in current state
		// The form validation should have already shown an error
		return nil

	case prefixConfirmedMsg:
		if m.config.Debug {
			spew.Fdump(m.config.DebugWriter, msg, "prefix confirmed")
		}

		if msg.confirmed {
			// User confirmed, save to database
			return m.transitionToOnboardingState(savingToDatabase, msg.prefix)
		}
		// User rejected, go back to selection
		return m.transitionToOnboardingState(selectingPrefix, "")

	case dbOperationCompleteMsg:
		if m.config.Debug {
			spew.Fdump(m.config.DebugWriter, msg, "db operation complete")
		}

		if msg.operation == "save" && msg.success {
			// Successfully saved network prefix, transition to operational
			return m.transitionToState(operational)
		}
		// Save failed, show error and stay in onboarding
		m.msg = fmt.Sprintf("Error saving to database: %v", msg.err)
		// Go back to prefix selection to try again
		return m.transitionToOnboardingState(selectingPrefix, "")
	}

	return nil
}

// handleOperationalMessage handles messages during operational state
func (m *mainModel) handleOperationalMessage(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case enterListViewMsg:
		if m.config.Debug {
			spew.Fdump(m.config.DebugWriter, msg, "entering list view")
		}
		return m.transitionToOperationalMode(listView)

	case enterDetailViewMsg:
		if m.config.Debug {
			spew.Fdump(m.config.DebugWriter, msg, "entering detail view")
		}
		m.currentRecordID = msg.recordID
		return m.transitionToOperationalMode(detailView)

	case enterCreateViewMsg:
		if m.config.Debug {
			spew.Fdump(m.config.DebugWriter, msg, "entering create view")
		}
		return m.transitionToOperationalMode(createView)

	case enterEditViewMsg:
		if m.config.Debug {
			spew.Fdump(m.config.DebugWriter, msg, "entering edit view")
		}
		m.currentRecordID = msg.recordID
		return m.transitionToOperationalMode(editView)

	case recordCreatedMsg:
		if m.config.Debug {
			spew.Fdump(m.config.DebugWriter, msg, "record created")
		}

		if msg.err == nil {
			// Success: show message and go back to list view
			m.msg = "Record created successfully"
			return func() tea.Msg { return enterListViewMsg{} }
		}
		// Handle error
		m.msg = fmt.Sprintf("Error creating record: %v", msg.err)
		return nil

	case recordUpdatedMsg:
		if m.config.Debug {
			spew.Fdump(m.config.DebugWriter, msg, "record updated")
		}

		if msg.err == nil {
			// Success: show message and go back to detail view
			m.msg = "Record updated successfully"
			return func() tea.Msg {
				return enterDetailViewMsg{recordID: msg.recordID}
			}
		}
		// Handle error
		m.msg = fmt.Sprintf("Error updating record: %v", msg.err)
		return nil

	case recordDeletedMsg:
		if m.config.Debug {
			spew.Fdump(m.config.DebugWriter, msg, "record deleted")
		}

		if msg.err == nil {
			// Success: show message and go back to list view
			m.msg = "Record deleted successfully"
			return func() tea.Msg { return enterListViewMsg{} }
		}
		// Handle error
		m.msg = fmt.Sprintf("Error deleting record: %v", msg.err)
		return nil

	case statusMsg:
		// Just display the status message
		m.msg = string(msg)
		return nil
	}

	return nil
}
