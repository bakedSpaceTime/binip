package app

import (
	"github.com/bakedSpaceTime/binip/libip/config"
	"github.com/bakedSpaceTime/binip/libip/db"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
)

type state uint

const (
	customPrefixStr       = "Custom Prefix"
	onboarding      state = iota
	operational
	quitting
)

type formState uint

const (
	formSelectPrefix formState = iota
	formCustomPrefix
	formConfirm
)

type mainModel struct {
	msg            string
	state          state
	formState      formState
	keys           keyMap
	help           help.Model
	lastKey        string
	width, height  int
	config         *config.Config
	firstWindowMsg bool
	db             *db.Db
	form           *huh.Form
	selectedPrefix string
	confirmed      bool
}

func New(c *config.Config) *mainModel {
	m := mainModel{
		state:          onboarding,
		keys:           keys,
		help:           help.New(),
		firstWindowMsg: true,
		config:         c,
		db:             db.New(c),
		formState:      formSelectPrefix,
	}

	m.form = m.onboardForm1()

	return &m
}

func (m *mainModel) Init() tea.Cmd {
	// start the timer, spinner, and api on program start
	return m.form.Init()
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.config.Debug {
		spew.Fdump(m.config.DebugWriter, msg, "main model message")
		m.msg = spew.Sdump(msg)
	}

	// Handle form updates during onboarding
	if m.state == onboarding {
		form, formCmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}

		// Check if form is completed
		if m.form.State == huh.StateCompleted {
			switch m.formState {
			case formSelectPrefix:
				if m.selectedPrefix == customPrefixStr {
					// User wants to enter custom prefix
					m.selectedPrefix = ""
					m.formState = formCustomPrefix
					m.form = m.onboardForm2()
					return m, m.form.Init()
				} else {
					// User selected a preset, go to confirmation
					m.confirmed = false
					m.formState = formConfirm
					m.form = m.onboardForm3()
					return m, m.form.Init()
				}

			case formCustomPrefix:
				// User entered custom prefix (already validated by form)
				// Go to confirmation
				m.confirmed = false
				m.formState = formConfirm
				m.form = m.onboardForm3()
				return m, m.form.Init()

			case formConfirm:
				// User confirmed or rejected the prefix
				if m.confirmed {
					// Save to database and transition to operational
					if err := m.db.SetNetwork(m.selectedPrefix); err != nil {
						if m.config.Debug {
							spew.Fdump(m.config.DebugWriter, err, "error saving network prefix")
						}
						m.msg = "Error: Failed to save network prefix"
						// Stay in onboarding state
						return m, tea.Batch(cmds...)
					}
					m.state = operational
					return m, nil
				} else {
					// User rejected, go back to form 1
					m.selectedPrefix = ""
					m.confirmed = false
					m.formState = formSelectPrefix
					m.form = m.onboardForm1()
					return m, m.form.Init()
				}
			}
		}

		cmds = append(cmds, formCmd)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width
		m.width = msg.Width
		m.height = msg.Height
		if m.firstWindowMsg {
			m.firstWindowMsg = false
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			m.state = quitting
			return m, tea.Quit
		}
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *mainModel) View() string {
	var body string
	switch m.state {
	case onboarding:
		body = m.onboardingView()
	default:
		body = m.db.String()
	}
	footer := m.footerView()
	return lipgloss.JoinVertical(lipgloss.Center, body, footer)

}
