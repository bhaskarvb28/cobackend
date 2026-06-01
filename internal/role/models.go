package role

import "cobackend/internal/shared/models"

type InvitableRoleOption struct {
	models.Role

	ScopeType string `json:"scope_type"`
}





