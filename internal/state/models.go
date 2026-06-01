package state

// State represents a state entity.
type State struct {

	// ID is the unique identifier of the state.
	ID int `json:"id"`

	// Name is the display name of the state.
	Name string `json:"name"`
}

// GetStatesQueryParams contains query parameters
// used for fetching states.
type GetStatesQueryParams struct {

	// Search filters states by name.
	Search string

	// Order defines the sorting order
	// of the returned states.
	Order string
}