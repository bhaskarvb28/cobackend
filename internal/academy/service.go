package academy

import "context"

func CreateAcademyService(
	ctx context.Context,
	input CreateAcademyInput,
) error {

	return CreateAcademyRepository(
		ctx,
		input,
	)
}

func GetAcademiesService(
	ctx context.Context,
	query GetAcademiesQuery,
) (PaginatedAcademies, error) {

	return GetAcademiesRepository(
		ctx,
		query,
	)
}