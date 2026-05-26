package role

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"cobackend/internal/db"
	"cobackend/internal/shared"
)

func GetRoleIDByName(
	ctx context.Context,
	roleName string,
) (string, error) {

	var roleID string

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT id
		FROM roles
		WHERE name = $1
		`,
		roleName,
	).Scan(&roleID)

	if err != nil {
		return "", err
	}

	return roleID, nil
}

func GetRoleNameByID(
	ctx context.Context,
	roleID string,
) (string, error) {
	var roleName string

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT name
		FROM roles
		WHERE id = $1
		`,
		roleID,
	).Scan(&roleName)

	if err != nil {
		return "", err
	}

	return roleName, nil
}

func GetRoleIDByCode(
	ctx context.Context,
	code string,
) (int, error) {

	query := `
		SELECT id
		FROM roles
		WHERE code = $1
	`

	var roleID int

	err := db.DB.QueryRow(
		ctx,
		query,
		code,
	).Scan(&roleID)

	if err != nil {

		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return 0, shared.ErrRoleNotFound
		}

		return 0, err
	}

	return roleID, nil
}