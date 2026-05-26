package academy

import (
	"context"
	"strings"

	"cobackend/internal/shared"
	"cobackend/internal/district"
)

func CreateAcademyService(
	ctx context.Context,
	input CreateAcademyInput,
) (*AcademyResponse, error) {

	// ----------------------------------------------------------
	// Normalize Input
	// ----------------------------------------------------------

	input.Name = strings.TrimSpace(input.Name)
	input.Address = strings.TrimSpace(input.Address)

	// ----------------------------------------------------------
	// Validate Input
	// ----------------------------------------------------------

	if input.Name == "" {
		return nil, shared.ErrAcademyNameRequired
	}

	if input.DistrictID <= 0 {
		return nil, shared.ErrInvalidDistrict
	}

	if input.Address == "" {
		return nil, shared.ErrAddressRequired
	}

	// ----------------------------------------------------------
	// Validate District Exists
	// ----------------------------------------------------------

	exists, err := district.CheckDistrictExists(
		ctx,
		input.DistrictID,
	)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, shared.ErrDistrictNotFound
	}

	// ----------------------------------------------------------
	// Create Academy
	// ----------------------------------------------------------

	academy, err := CreateAcademyRepository(
		ctx,
		input,
	)

	if err != nil {
		return nil, err
	}

	return academy, nil
}

// func GetAcademiesService(
// 	ctx context.Context,
// 	query GetAcademiesQuery,
// ) (PaginatedAcademies, error) {

// 	return GetAcademiesRepository(
// 		ctx,
// 		query,
// 	)
// }