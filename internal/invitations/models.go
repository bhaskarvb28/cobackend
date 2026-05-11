package invitations

import (
	"time"
)

type Invitation struct {
	ID               string
	Email            string
	RoleID           string
	Token            string
	AssignedStateID  string
	Status           string
	ExpiresAt        time.Time
}