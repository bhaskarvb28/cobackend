package role

import (
	"context"
	"errors"

	"cobackend/internal/db"
	"cobackend/internal/shared"
	"cobackend/internal/shared/models"

	"github.com/jackc/pgx/v5"
)

func GetRolesByCodesRepository(
	ctx context.Context,
	roleCodes []string,
) ([]models.Role, error) {

	query := `
		SELECT
			id,
			code,
			display_name
		FROM roles
		WHERE code = ANY($1)
		ORDER BY id ASC
	`

	rows, err := db.DB.Query(
		ctx,
		query,
		roleCodes,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	roles := []models.Role{}

	for rows.Next() {

		var role models.Role

		err := rows.Scan(
			&role.ID,
			&role.Code,
			&role.DisplayName,
		)

		if err != nil {
			return nil, err
		}

		roles = append(
			roles,
			role,
		)
	}

	return roles, nil
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

// func GetRoleIDByName(
// 	ctx context.Context,
// 	roleName string,
// ) (string, error) {

// 	var roleID string

// 	err := db.DB.QueryRow(
// 		ctx,
// 		`
// 		SELECT id
// 		FROM roles
// 		WHERE name = $1
// 		`,
// 		roleName,
// 	).Scan(&roleID)

// 	if err != nil {
// 		return "", err
// 	}

// 	return roleID, nil
// }

// func GetRoleNameByID(
// 	ctx context.Context,
// 	roleID string,
// ) (string, error) {
// 	var roleName string

// 	err := db.DB.QueryRow(
// 		ctx,
// 		`
// 		SELECT name
// 		FROM roles
// 		WHERE id = $1
// 		`,
// 		roleID,
// 	).Scan(&roleName)

// 	if err != nil {
// 		return "", err
// 	}

// 	return roleName, nil
// }



func GetRolesRepository(
	ctx context.Context,
) ([]models.Role, error) {

	query := `
		SELECT
			id,
			code,
			display_name
		FROM roles
		ORDER BY id ASC
	`

	rows, err := db.DB.Query(
		ctx,
		query,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	roles := []models.Role{}

	for rows.Next() {

		var role models.Role

		err := rows.Scan(
			&role.ID,
			&role.Code,
			&role.DisplayName,
		)

		if err != nil {
			return nil, err
		}

		roles = append(
			roles,
			role,
		)
	}

	return roles, nil
}