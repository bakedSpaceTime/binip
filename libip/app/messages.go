package app

// === State Transition Messages ===

// stateTransitionMsg indicates a transition between main states
type stateTransitionMsg struct {
	fromState state
	toState   state
}

// === Onboarding Messages ===

// prefixSelectedMsg is sent when user selects a prefix from the initial form
type prefixSelectedMsg struct {
	prefix   string
	isCustom bool
}

// customPrefixEnteredMsg is sent when user enters a custom prefix
type customPrefixEnteredMsg struct {
	prefix string
	valid  bool
}

// prefixConfirmedMsg is sent when user confirms or rejects the prefix
type prefixConfirmedMsg struct {
	confirmed bool
	prefix    string
}

// === Async Operation Result Messages ===

// dbOperationCompleteMsg is sent when a database operation completes
type dbOperationCompleteMsg struct {
	operation string      // "save", "update", "delete"
	success   bool        // whether operation succeeded
	err       error       // error if operation failed
	data      interface{} // optional data returned from operation
}

// === Operational CRUD Messages ===

// enterListViewMsg requests transition to list view
type enterListViewMsg struct{}

// enterDetailViewMsg requests transition to detail view for a specific record
type enterDetailViewMsg struct {
	recordID string
}

// enterCreateViewMsg requests transition to create view
type enterCreateViewMsg struct{}

// enterEditViewMsg requests transition to edit view for a specific record
type enterEditViewMsg struct {
	recordID string
}

// recordCreatedMsg is sent when a record is created
type recordCreatedMsg struct {
	record interface{} // the created record (replace with your record type)
	err    error
}

// recordUpdatedMsg is sent when a record is updated
type recordUpdatedMsg struct {
	recordID string
	err      error
}

// recordDeletedMsg is sent when a record is deleted
type recordDeletedMsg struct {
	recordID string
	err      error
}

// === Error/Status Messages ===

// errorMsg represents an error with context
type errorMsg struct {
	context string
	err     error
}

// statusMsg represents a status message to display to user
type statusMsg string
