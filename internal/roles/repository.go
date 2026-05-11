package roles

import (
	"cobackend/internal/db"
	"context"
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