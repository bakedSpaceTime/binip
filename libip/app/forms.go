package app

import (
	"fmt"
	"net/netip"

	"github.com/bakedSpaceTime/binip/libip/record"
	"github.com/charmbracelet/huh"
)

func (m *mainModel) onboardForm1() *huh.Form {
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
				Value(&m.selectedPrefix),
		),
	)
}

func (m *mainModel) onboardForm2() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter custom network prefix").
				Placeholder("e.g., 192.168.1.0/24").
				Description("Must be in CIDR notation").
				Value(&m.selectedPrefix).
				Validate(func(s string) error {
					if _, err := netip.ParsePrefix(s); err != nil {
						return fmt.Errorf("invalid CIDR notation - use format like 192.168.1.0/24")
					}
					return nil
				}),
		),
	)
}

func (m *mainModel) onboardForm3() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Use this network prefix?\n\n  %s", m.selectedPrefix)).
				Affirmative("Yes").
				Negative("No, choose again").
				Value(&m.confirmed),
		),
	)
}
