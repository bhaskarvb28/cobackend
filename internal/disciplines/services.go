package disciplines

import (
	"context"
)

func GetDisciplinesService(ctx context.Context) ([]Disciplines, error) {
	return GetDisciplinesRepository(ctx)
}