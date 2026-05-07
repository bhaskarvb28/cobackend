package states

import "context"

func GetStatesService(ctx context.Context) ([]State, error) {
	return GetStatesRepository(ctx)
}