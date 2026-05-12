package districts

import (
	"context"
)

func GetDistrictsService(ctx context.Context) ([]District, error) {
	return GetDistrictsRepository(ctx)
}

func GetDistrictsByStateIdService(ctx context.Context, stateID int, queryParams GetDistrictQueryParams) ([]DistrictResponse, error) {
	return GetDistrictsByStateIDRepository(ctx, stateID, queryParams)
}