package auth

import (
	"cobackend/internal/shared/models"
)

// LoginInput contains user credentials
// required for authentication
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthUser contains authentication-related
// user data loaded during login validation.
type AuthUser struct {
	ID            string
	FirstName     string
	LastName      *string
	Email         string
	ContactNumber string
	PasswordHash  string
	Role          models.Role
}

// LoginResponse contains authentication tokens
// and authenticated user information.
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponse contains authenticated user data.
// UserResponse contains authenticated user data.
type UserResponse struct {
	ID            string     `json:"id"`
	FirstName     string     `json:"first_name"`
	LastName      *string    `json:"last_name,omitempty"`
	Email         string     `json:"email"`
	ContactNumber string     `json:"contact_number"`
	Role          models.Role `json:"role"`
}

// ----------------------------------------------------------------------------------------------------------

type CreateUserInput struct {
	FirstName     string
	LastName      string
	Email         string
	PasswordHash  string
	ContactNumber string
	Role          string
}

// type AcceptInvitationInput struct {
// 	Token         string `json:"token"`
// 	FirstName     string `json:"first_name"`
// 	LastName      string `json:"last_name"`
// 	Password      string `json:"password"`
// 	ContactNumber string `json:"contact_number"`

// 	DPDPConsent bool `json:"dpdp_consent,omitempty"`

// 	CoachCode                string `json:"coach_code,omitempty"`
// 	CoachingCertificateProof string `json:"coaching_certificate_proof,omitempty"`

// 	GSTIN              string `json:"gstin,omitempty"`
// 	RegistrationProof  string `json:"registration_proof,omitempty"`

// 	CoachingCredentialsProof string `json:"coaching_credentials_proof,omitempty"`

// 	//------------------------------------------------
// 	// Player Fields
// 	//------------------------------------------------

// 	DateOfBirth string `json:"date_of_birth,omitempty"`

// 	Gender string `json:"gender,omitempty"`

// 	ParentGuardianName string `json:"parent_guardian_name,omitempty"`

// 	ParentGuardianContact string `json:"parent_guardian_contact,omitempty"`

// 	AlternativeContact string `json:"alternative_contact,omitempty"`

// 	ParentalConsent bool `json:"parental_consent,omitempty"`

// }