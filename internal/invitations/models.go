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
	Used             bool
	ExpiresAt        time.Time
}