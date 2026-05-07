package districts

import "context"

func GetDistrictsService(ctx context.Context) ([]District, error) {
	return GetDistrictsRepository(ctx)
}