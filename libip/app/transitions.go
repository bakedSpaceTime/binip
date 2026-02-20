package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// canTransition validates if a main state transition is allowed
func (m *mainModel) canTransition(from, to state) error {
	switch {
	case from == quitting:
		return fmt.Errorf("cannot transition from quitting state")
	case to == operational && m.db == nil:
		return fmt.Errorf("database required before operational state")
	}
	return nil
}

// validateOnboardingTransition validates if an onboarding sub-state transition is allowed
func (m *mainModel) validateOnboardingTransition(to onboardingState, prefix string) error {
	switch to {
	case confirmingPrefix:
		if prefix == "" {
			return fmt.Errorf("prefix must be selected or entered")
		}
	case savingToDatabase:
		if prefix == "" {
			return fmt.Errorf("prefix required before saving")
		}
	}
	return nil
}

// transitionToState transitions to a new main state with validation
func (m *mainModel) transitionToState(newState state) tea.Cmd {
	if err := m.canTransition(m.state, newState); err != nil {
		return func() tea.Msg {
			return errorMsg{context: "state transition", err: err}
		}
	}

	oldState := m.state
	m.state = newState

	// State entry actions
	var cmd tea.Cmd
	switch newState {
	case operational:
		m.operationalMode = listView
		cmd = m.loadOperationalData()
	}

	// Return a batch of the state transition message and any entry action command
	return tea.Batch(
		func() tea.Msg {
			return stateTransitionMsg{fromState: oldState, toState: newState}
		},
		cmd,
	)
}

// transitionToOnboardingState transitions to a new onboarding sub-state
func (m *mainModel) transitionToOnboardingState(newState onboardingState, prefix string) tea.Cmd {
	if err := m.validateOnboardingTransition(newState, prefix); err != nil {
		return func() tea.Msg {
			return errorMsg{context: "onboarding transition", err: err}
		}
	}

	m.onboardingState = newState

	// Initialize the appropriate form/view for the new state
	switch newState {
	case selectingPrefix:
		m.form = m.prefixSelectForm()
		return m.form.Init()

	case enteringCustomPrefix:
		m.form = m.customPrefixForm()
		return m.form.Init()

	case confirmingPrefix:
		// Store the prefix being confirmed so it's available when form completes
		m.prefixBeingConfirmed = prefix
		m.form = m.confirmPrefixForm(prefix)
		return m.form.Init()

	case savingToDatabase:
		return m.saveNetworkPrefix(prefix)
	}

	return nil
}

// transitionToOperationalMode transitions to a different operational mode
func (m *mainModel) transitionToOperationalMode(newMode operationalMode) tea.Cmd {
	m.operationalMode = newMode

	// Mode entry actions
	switch newMode {
	case listView:
		return m.loadRecordList()

	case detailView:
		// Requires a record ID to be set before calling this
		if m.currentRecordID != "" {
			return m.loadRecordDetail(m.currentRecordID)
		}

	case createView:
		m.form = m.createRecordForm()
		return m.form.Init()

	case editView:
		// Requires a record ID to be set before calling this
		if m.currentRecordID != "" {
			m.form = m.editRecordForm(m.currentRecordID)
			return m.form.Init()
		}

	case deleteConfirmView:
		// Requires a record ID to be set before calling this
		if m.currentRecordID != "" {
			m.form = m.deleteConfirmForm(m.currentRecordID)
			return m.form.Init()
		}
	}

	return nil
}
