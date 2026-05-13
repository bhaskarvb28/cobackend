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