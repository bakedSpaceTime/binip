package app

import (
	"fmt"

	"github.com/bakedSpaceTime/binip/libip/config"
	"github.com/bakedSpaceTime/binip/libip/db"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
)

const (
	customPrefixStr = "Custom Prefix"
)

type mainModel struct {
	// Main state
	state           state
	onboardingState onboardingState
	operationalMode operationalMode

	// Temporary form binding fields (not persisted - data flows through messages)
	formPrefix           string // Temporary: binds to prefix selection/input forms
	formConfirmed        bool   // Temporary: binds to confirmation forms
	prefixBeingConfirmed string // Temporary: holds prefix during confirmation flow

	// Operational data
	currentRecordID string        // Currently selected record ID
	records         []interface{} // List of records (replace with your record type)

	// UI components
	form   *huh.Form
	keys   keyMap
	help   help.Model
	width  int
	height int
	msg    string // Status or error message to display

	// Dependencies
	config *config.Config
	db     *db.Db

	// Internal
	firstWindowMsg bool
	lastKey        string
}

func New(c *config.Config) *mainModel {
	m := mainModel{
		state:           onboarding,
		onboardingState: selectingPrefix,
		keys:            keys,
		help:            help.New(),
		config:          c,
		db:              db.New(c),
		firstWindowMsg:  true,
	}

	m.form = m.prefixSelectForm()
	return &m
}

func (m *mainModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Debug logging
	if m.config.Debug {
		spew.Fdump(m.config.DebugWriter, msg, "message received")
		m.msg = spew.Sdump(msg)
	}

	// Global message handling (window size, keyboard shortcuts, errors)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.help.Width = msg.Width
		m.firstWindowMsg = false

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			m.state = quitting
			return m, tea.Quit
		}

	case errorMsg:
		// Global error handling
		m.msg = fmt.Sprintf("Error [%s]: %v", msg.context, msg.err)
		if m.config.Debug {
			spew.Fdump(m.config.DebugWriter, msg.err, "error occurred")
		}
		return m, nil

	case stateTransitionMsg:
		// Log state transitions
		if m.config.Debug {
			spew.Fdump(m.config.DebugWriter, msg, "state transition")
		}
	}

	// Update form if one is active and not completed
	if m.form != nil && m.form.State != huh.StateCompleted {
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}
		cmds = append(cmds, cmd)

		// Check for form completion and handle it
		if m.form.State == huh.StateCompleted {
			cmds = append(cmds, m.handleFormCompletion())
		}
	}

	// State-specific message handling
	cmds = append(cmds, m.handleMessage(msg))

	return m, tea.Batch(cmds...)
}

func (m *mainModel) View() string {
	var body string

	switch m.state {
	case onboarding:
		body = m.onboardingView()
	case operational:
		body = m.operationalView()
	default:
		body = ""
	}

	footer := m.footerView()
	return lipgloss.JoinVertical(lipgloss.Center, body, footer)
}

// handleFormCompletion handles form completion by returning appropriate messages
// based on the current state
func (m *mainModel) handleFormCompletion() tea.Cmd {
	switch m.state {
	case onboarding:
		return m.handleOnboardingFormCompletion()
	case operational:
		return m.handleOperationalFormCompletion()
	}
	return nil
}

// handleOnboardingFormCompletion handles form completion during onboarding
func (m *mainModel) handleOnboardingFormCompletion() tea.Cmd {
	switch m.onboardingState {
	case selectingPrefix:
		// User selected a prefix (preset or custom option)
		return func() tea.Msg {
			return prefixSelectedMsg{
				prefix:   m.formPrefix,
				isCustom: m.formPrefix == customPrefixStr,
			}
		}

	case enteringCustomPrefix:
		// User entered a custom prefix (already validated by form)
		return func() tea.Msg {
			return customPrefixEnteredMsg{
				prefix: m.formPrefix,
				valid:  true, // Form validation ensures this
			}
		}

	case confirmingPrefix:
		// User confirmed or rejected the prefix
		// Use the prefix that was stored when we transitioned to this state
		return func() tea.Msg {
			return prefixConfirmedMsg{
				confirmed: m.formConfirmed,
				prefix:    m.prefixBeingConfirmed,
			}
		}
	}

	return nil
}

// handleOperationalFormCompletion handles form completion during operational state
func (m *mainModel) handleOperationalFormCompletion() tea.Cmd {
	switch m.operationalMode {
	case createView:
		// TODO: Extract form data and create record
		return func() tea.Msg {
			return recordCreatedMsg{record: nil, err: nil}
		}

	case editView:
		// TODO: Extract form data and update record
		return func() tea.Msg {
			return recordUpdatedMsg{recordID: m.currentRecordID, err: nil}
		}

	case deleteConfirmView:
		if m.formConfirmed {
			return m.deleteRecord(m.currentRecordID)
		}
		// User cancelled, go back to detail view
		return func() tea.Msg {
			return enterDetailViewMsg{recordID: m.currentRecordID}
		}
	}

	return nil
}
