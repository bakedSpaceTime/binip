package app

// state represents the main application state
type state uint

const (
	onboarding state = iota
	operational
	quitting
)

func (s state) String() string {
	switch s {
	case onboarding:
		return "onboarding"
	case operational:
		return "operational"
	case quitting:
		return "quitting"
	default:
		return "unknown"
	}
}

// onboardingState represents sub-states during onboarding flow
type onboardingState uint

const (
	selectingPrefix onboardingState = iota
	enteringCustomPrefix
	confirmingPrefix
	savingToDatabase // Explicit "in progress" state for async save
)

func (os onboardingState) String() string {
	switch os {
	case selectingPrefix:
		return "selecting prefix"
	case enteringCustomPrefix:
		return "entering custom prefix"
	case confirmingPrefix:
		return "confirming prefix"
	case savingToDatabase:
		return "saving to database"
	default:
		return "unknown"
	}
}

// operationalMode represents different views/modes in operational state
// This is extensible for future CRUD operations
type operationalMode uint

const (
	listView operationalMode = iota
	detailView
	createView
	editView
	deleteConfirmView
	// Easy to add more modes as UI design evolves
)

func (om operationalMode) String() string {
	switch om {
	case listView:
		return "list view"
	case detailView:
		return "detail view"
	case createView:
		return "create view"
	case editView:
		return "edit view"
	case deleteConfirmView:
		return "delete confirm view"
	default:
		return "unknown"
	}
}
