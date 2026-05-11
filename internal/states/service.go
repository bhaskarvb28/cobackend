package states

import "context"

func GetStatesService(ctx context.Context, query GetStatesQueryParams) ([]State, error) {
	return GetStatesRepository(ctx, query)
}