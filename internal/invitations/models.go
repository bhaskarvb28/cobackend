package invitations

import (
	"time"
)

type Invitation struct {
	ID    string
	Email string

	RoleID   string
	RoleName string

	Token  string
	Status string

	StateID    *string
	DistrictID *string
	AcademyID  *string

	ExpiresAt time.Time
}