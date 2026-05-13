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

	StateID    *int
	DistrictID *int
	AcademyID  *int

	ExpiresAt time.Time
}