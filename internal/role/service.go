package role

import (
	"context"

	"cobackend/internal/shared"
)

type invitableRoleConfig struct {
	Code      string
	ScopeType string
}

var invitableRolesMap = map[string][]invitableRoleConfig{

	// ----------------------------------------------------------
	// Super Admin
	// ----------------------------------------------------------

	"super_admin": {
		{
			Code:      "state_admin",
			ScopeType: "state",
		},
	},

	// ----------------------------------------------------------
	// State Admin
	// ----------------------------------------------------------

	"state_admin": {
		{
			Code:      "district_admin",
			ScopeType: "district",
		},
		{
			Code:      "district_coach",
			ScopeType: "district",
		},
	},

	// ----------------------------------------------------------
	// District Admin
	// ----------------------------------------------------------

	"district_admin": {
		{
			Code:      "academy_admin",
			ScopeType: "academy",
		},
	},

	// ----------------------------------------------------------
	// Academy Admin
	// ----------------------------------------------------------

	"academy_admin": {
		{
			Code:      "academy_coach",
			ScopeType: "academy",
		},
		{
			Code:      "player",
			ScopeType: "academy",
		},
	},

	// ----------------------------------------------------------
	// No Invite Permissions
	// ----------------------------------------------------------

	"district_coach": {},
	"academy_coach": {},
	"player":         {},
}

func GetInvitableRoleOptionsService(
	ctx context.Context,
	roleCode string,
) ([]InvitableRoleOption, error) {

	configs, ok := invitableRolesMap[roleCode]

	if !ok {
		return nil, shared.ErrInvalidRole
	}

	// ----------------------------------------------------------
	// Extract Role Codes
	// ----------------------------------------------------------

	roleCodes := []string{}

	for _, config := range configs {
		roleCodes = append(
			roleCodes,
			config.Code,
		)
	}

	// ----------------------------------------------------------
	// Fetch Roles From DB
	// ----------------------------------------------------------

	roles, err := GetRolesByCodesRepository(
		ctx,
		roleCodes,
	)

	if err != nil {
		return nil, err
	}

	// ----------------------------------------------------------
	// Create Lookup Map
	// ----------------------------------------------------------

	roleMap := map[string]InvitableRoleOption{}

	for _, role := range roles {

		roleMap[role.Code] = InvitableRoleOption{
			Role: role,
		}
	}

	// ----------------------------------------------------------
	// Attach Scope Types
	// ----------------------------------------------------------

	result := []InvitableRoleOption{}

	for _, config := range configs {

		roleOption := roleMap[config.Code]

		roleOption.ScopeType =
			config.ScopeType

		result = append(
			result,
			roleOption,
		)
	}

	return result, nil
}
