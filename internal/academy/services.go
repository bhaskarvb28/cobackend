package academy

import (
	"context"
	"strings"

	"cobackend/internal/districtAdmin"
	"cobackend/internal/shared"
)

func CreateAcademyService(
	ctx context.Context,
	userID string,
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

	if input.Address == "" {
		return nil, shared.ErrAddressRequired
	}

	// ----------------------------------------------------------
	// Get District Admin
	// ----------------------------------------------------------

	districtAdminRegion, err := districtAdmin.GetDistrictAdminRegion(
		ctx,
		userID,
	)

	if err != nil {
		return nil, err
	}

	// ----------------------------------------------------------
	// Create Academy
	// ----------------------------------------------------------

	academy, err := CreateAcademyRepository(
		ctx,
		districtAdminRegion.DistrictID,
		input,
	)

	if err != nil {
		return nil, err
	}

	return academy, nil
}

// GetAcademiesService retrieves a paginated list of academies
// based on state, district, search, and sorting criteria.
func GetAcademiesService(
	ctx context.Context,
	query GetAcademiesQuery,
) (PaginatedAcademies, error) {

	// ----------------------------------------------------------
	// Fetch Academies
	// ----------------------------------------------------------
	// Query repository layer directly to fetch rows and pagination total counts
	return GetAcademiesRepository(
		ctx,
		query,
	)
}