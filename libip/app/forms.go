package app

import (
	"fmt"
	"net/netip"

	"github.com/bakedSpaceTime/binip/libip/record"
	"github.com/charmbracelet/huh"
)

// === Onboarding Forms ===

// prefixSelectForm shows the initial prefix selection form
func (m *mainModel) prefixSelectForm() *huh.Form {
	options := make([]huh.Option[string], len(record.PrivateRanges)+1)
	for i, pr := range record.PrivateRanges {
		options[i] = huh.NewOption(pr.String(), pr.String())
	}
	options[len(record.PrivateRanges)] = huh.NewOption(customPrefixStr, customPrefixStr)

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select IP Range").
				Options(options...).
				Value(&m.formPrefix), // Bind to temporary form field
		),
	)
}

// customPrefixForm allows user to enter a custom prefix
func (m *mainModel) customPrefixForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter custom network prefix").
				Placeholder("e.g., 192.168.1.0/24").
				Description("Must be in CIDR notation").
				Value(&m.formPrefix). // Bind to temporary form field
				Validate(func(s string) error {
					if _, err := netip.ParsePrefix(s); err != nil {
						return fmt.Errorf("invalid CIDR notation - use format like 192.168.1.0/24")
					}
					return nil
				}),
		),
	)
}

// confirmPrefixForm shows confirmation for the selected prefix
func (m *mainModel) confirmPrefixForm(prefix string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Use this network prefix?\n\n  %s", prefix)).
				Affirmative("Yes").
				Negative("No, choose again").
				Value(&m.formConfirmed), // Bind to temporary form field
		),
	)
}

// === CRUD Forms (Templates for future implementation) ===

// createRecordForm creates a form for creating a new record
func (m *mainModel) createRecordForm() *huh.Form {
	// TODO: Implement when CRUD UI is designed
	return huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Create Record").
				Description("Form to be designed"),
		),
	)
}

// editRecordForm creates a form for editing an existing record
func (m *mainModel) editRecordForm(recordID string) *huh.Form {
	// TODO: Implement when CRUD UI is designed
	// Load the record and populate form fields
	return huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Edit Record").
				Description("Form to be designed for record: " + recordID),
		),
	)
}

// deleteConfirmForm creates a confirmation form for deleting a record
func (m *mainModel) deleteConfirmForm(recordID string) *huh.Form {
	// TODO: Implement when CRUD UI is designed
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Delete record %s?", recordID)).
				Affirmative("Yes, delete").
				Negative("Cancel").
				Value(&m.formConfirmed),
		),
	)
}
