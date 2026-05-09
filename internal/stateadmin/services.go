package stateadmin

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"cobackend/internal/auth"
	"cobackend/internal/shared"
	"cobackend/internal/states"
)

func CreateStateAdminService(
	ctx context.Context,
	input CreateStateAdminInput,
) error {

	// check if email already exists
	exists, err := auth.CheckEmailExists(ctx, input.Email)
	if err != nil {
		return err
	}

	if exists {
		return shared.ErrEmailExists
	}

	// validate assigned state exists
	stateExists, err := states.CheckStateExists(
		ctx,
		input.AssignedState,
	)

	if err != nil {
		return err
	}

	if !stateExists {
		return shared.ErrInvalidState
	}

	// hash password
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

func GetStateAdminsService(ctx context.Context, query GetStateAdminsQuery) ([]StateAdminResponse, error) {
	return GetStateAdminsRepository(ctx, query)
}


func UpdateAssignedStateService(
	ctx context.Context,
	id string,
	input UpdateAssignedStateInput,
) error {

	// check state admin exists
	exists, err := CheckStateAdminExists(
		ctx,
		id,
	)

	if err != nil {
		return err
	}

	if !exists {
		return shared.ErrStateAdminNotFound
	}

	// validate state exists
	stateExists, err := states.CheckStateExists(
		ctx,
		input.AssignedState,
	)

	if err != nil {
		return err
	}

	if !stateExists {
		return shared.ErrInvalidState
	}

	return UpdateAssignedStateRepository(
		ctx,
		id,
		input,
	)
}


func DeleteStateAdminService(
	ctx context.Context,
	id string,
) error {

	// check state admin exists
	exists, err := CheckStateAdminExists(
		ctx,
		id,
	)

	if err != nil {
		return err
	}

	if !exists {
		return shared.ErrStateAdminNotFound
	}

	return DeleteStateAdminRepository(
		ctx,
		id,
	)
}