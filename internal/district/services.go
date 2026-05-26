package district

import "context"

// GetDistrictsService fetches all districts.
func GetDistrictsService(
	ctx context.Context,
) ([]District, error) {

	return GetDistrictsRepository(
		ctx,
	)
}

// GetDistrictsByStateIdService fetches all districts
// belonging to a specific state based on
// the provided query parameters.
func GetDistrictsByStateIdService(
	ctx context.Context,
	stateID int,
	queryParams GetDistrictQueryParams,
) ([]DistrictResponse, error) {

	return GetDistrictsByStateIDRepository(
		ctx,
		stateID,
		queryParams,
	)
}