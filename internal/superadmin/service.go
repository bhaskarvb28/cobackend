package superadmin

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

func CreateStateAdminService(
	ctx context.Context,
	input CreateStateAdminInput,
) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	return CreateStateAdminRepository(
		ctx,
		input,
		string(hashedPassword),
	)
}

func UpdateAssignedStateService(
	ctx context.Context,
	id string,
	input UpdateAssignedStateInput,
) error {

	return UpdateAssignedStateRepository(
		ctx,
		id,
		input,
	)
}

func GetStateAdminsService(ctx context.Context, query GetStateAdminsQuery) ([]StateAdmin, error) {
	return GetStateAdminsRepository(ctx, query)
}

func DeleteStateAdminService(ctx context.Context, id string) error {
	return DeleteStateAdminRepository(ctx, id)
}