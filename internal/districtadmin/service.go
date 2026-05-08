package districtadmin

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

func CreateDistrictAdminService(ctx context.Context, input CreateDistrictAdminInput) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	return CreateDistrictAdminRepository(ctx, input, string(hashedPassword))
}

func GetDistrictAdminsService(ctx context.Context, query GetDistrictAdminsQuery) ([]DistrictAdmin, error) {
	return GetDistrictAdminsRepository(ctx, query)
}

func UpdateDistrictAdminService(ctx context.Context, id string, input UpdateDistrictAdminInput) error {
	return UpdateDistrictAdminRepository(ctx, id, input)
}

func DeleteDistrictAdminService(ctx context.Context, id string) error {
	return DeleteDistrictAdminRepository(ctx, id)
}
