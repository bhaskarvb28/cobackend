package invitation

import "time"

// InvitationQueryParams contains
// query parameters for fetching invitations.
type InvitationsQueryParams struct {

	// Pagination
	Page  int `query:"page"`
	Limit int `query:"limit"`

	// Filters
	Search string `query:"search"`
	Status string `query:"status"`
	Role   string `query:"role"`

	// Sorting
	SortBy string `query:"sort_by"`
	Order  string `query:"order"`
}

type PaginatedInvitations struct {
	Items       []InvitationResponse `json:"items"`
	Page        int          `json:"page"`
	Limit       int          `json:"limit"`
	Total       int          `json:"total"`
	TotalPages  int          `json:"total_pages"`
	HasNext     bool         `json:"has_next"`
	HasPrevious bool         `json:"has_previous"`
}

// Invitation represents an invitation entity.
type Invitation struct {

	// ID is the unique identifier
	// of the invitation.
	ID int64 `json:"id"`

	// Email is the invited user's email.
	Email string `json:"email"`

	// Role is the assigned role code.
	//
	// Examples:
	//	- super_admin
	//	- state_admin
	//	- district_admin
	//	- academy_admin
	//	- coach
	//	- player
	Role string `json:"role"`

	// InvitedBy is the user ID of
	// the user who created the invitation.
	InvitedBy string `json:"invited_by"`

	// ScopeType defines the type
	// of resource scope assigned
	// through the invitation.
	//
	// Examples:
	//	- state
	//	- district
	//	- academy
	ScopeType *string `json:"scope_type,omitempty"`

	// ScopeID identifies the specific
	// scoped resource.
	//
	// Examples:
	//	- "1"
	//	- "15"
	//	- "uuid"
	ScopeID *string `json:"scope_id,omitempty"`

	// ExpiresAt defines when the
	// invitation expires.
	ExpiresAt time.Time `json:"expires_at"`

	// Status represents the current
	// invitation status.
	//
	// Possible values:
	//	- pending
	//	- accepted
	//	- expired
	//	- revoked
	Status string `json:"status"`

	// AcceptedAt stores when
	// the invitation was accepted.
	AcceptedAt *time.Time `json:"accepted_at,omitempty"`

	// UsedBy stores the user ID
	// of the user who accepted
	// the invitation.
	UsedBy *string `json:"used_by,omitempty"`

	// CreatedAt stores the
	// invitation creation timestamp.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt stores the last
	// update timestamp.
	UpdatedAt time.Time `json:"updated_at"`
}

type InvitationPermissionCheck struct {
	ID        int64
	InvitedBy string
	RoleCode  string
	Status    string
	ExpiresAt time.Time
}

type RoleResponse struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

type OrganizationResponse struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// InvitationResponse represents
// public invitation response data.
type InvitationResponse struct {

	ID int64 `json:"id"`

	Email string `json:"email"`

	Role RoleResponse `json:"role"`

	RoleCode string `json:"role_code"`

	Organization *OrganizationResponse `json:"organization,omitempty"`

	ScopeType *string `json:"scope_type,omitempty"`

	ScopeID *string `json:"scope_id,omitempty"`

	Status string `json:"status"`

	CreatedBy UserSummary `json:"created_by"`

	InvitedBy string `json:"invited_by"`

	ExpiresAt time.Time `json:"expires_at"`

	AcceptedAt *time.Time `json:"accepted_at,omitempty"`

	UsedBy *string `json:"used_by,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

// CreateInvitationInput represents
// the request body for creating
// an invitation.
type CreateInvitationInput struct {

	// Name is the invited user's name.
	Name string `json:"name"`

	// Email is the invited user's email.
	Email string `json:"email"`

	// Role is the target role code.
	//
	// Examples:
	//	- state_admin
	//	- district_admin
	//	- academy_admin
	Role string `json:"role"`

	// ScopeType defines the scope type.
	//
	// Examples:
	//	- state
	//	- district
	//	- academy
	ScopeType string `json:"scope_type"`

	// ScopeID defines the scoped entity ID.
	ScopeID string `json:"scope_id"`
}

//------------------------------------------------------------------------------------------------
//BORDER
//------------------------------------------------------------------------------------------------








// AcceptInvitationInput represents the request
// body for accepting an invitation.
type AcceptInvitationInput struct {

	// Token is the invitation token.
	Token string `json:"token"`

	// FirstName is the user's first name.
	FirstName string `json:"first_name"`

	// LastName is the user's last name.
	LastName string `json:"last_name"`

	// Password is the account password.
	Password string `json:"password"`

	// ContactNumber is the user's contact number.
	ContactNumber string `json:"contact_number"`
}