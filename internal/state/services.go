package state

import "context"

// GetStatesService fetches all states
// based on the provided query parameters.
func GetStatesService(
	ctx context.Context,
	query GetStatesQueryParams,
) ([]State, error) {

	return GetStatesRepository(
		ctx,
		query,
	)
}