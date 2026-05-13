package auth

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Internal auth/database model
type AuthUser struct {
	ID           string
	Email        string
	PasswordHash string
	RoleID       string
	Role         string
}

// Login API response
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// Public user response
type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type AcceptInvitationInput struct {
	Token         string `json:"token"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Password      string `json:"password"`
	ContactNumber string `json:"contact_number"`

	DPDPConsent bool `json:"dpdp_consent,omitempty"`
	CoachCode                string `json:"coach_code"`
	CoachingCertificateProof string `json:"coaching_certificate_proof"`


	// Specialization string `json:"specialization,omitempty"`

	// DateOfBirth string `json:"date_of_birth,omitempty"`

	// Gender string `json:"gender,omitempty"`
}
