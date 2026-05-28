package shared

import "errors"

var (
	
	// Midddleware
	// AuthMiddleware
	ErrMissingAuthorizationHeader = errors.New("missing authorization header")
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
	ErrInvalidToken = errors.New("invalid token")

	// RequiredRoleMiddleware
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden = errors.New("forbidden")
	ErrInvalidScope = errors.New("invalid scope")
	ErrRoleNotAuthorized = errors.New("role is not authorized for this action")


	// Generic Errors
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrInternalServerError = errors.New("internal server error")
	ErrFirstNameRequired = errors.New("first name is required")
	ErrPasswordRequired = errors.New("password is required")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrDPDPConsentRequired = errors.New("dpdp consent is required")

	// Authentication
	ErrInvalidCredentials = errors.New("invalid email or password")


	// Invitation
	ErrInvitationAlreadyExists = errors.New("pending invitation already exists")
	ErrInvitationExpired = errors.New("invitation expired")
	ErrInvitationNotFound = errors.New("invitation not found")
	ErrInvitationRevoked = errors.New("invitation has been revoked")
	ErrInvitationAlreadyAccepted = errors.New("invitation has already been acccepted")
	ErrInvitationAlreadyRevoked = errors.New("invitation has already been revoked")
	ErrTokenRequired = errors.New("token is required")
	ErrInvalidInvitationScope = errors.New("invalid invitation scope")

	// Districts
	ErrDistrictNotFound = errors.New("district not found")

	// Academy
	ErrAcademyNotFound = errors.New("academy not found")
	ErrAcademyNameRequired = errors.New("academy name required")
	ErrAddressRequired = errors.New("academy address required")

	// Academy Building
	ErrAcademyBuildingNameRequired = errors.New("academy building name required")
	ErrUnauthorizedBuildingAccess = errors.New("unauthorized building access")
	ErrInvalidBuildingID = errors.New("invalid building id")

	// Discipline
	ErrInvalidDisciplineID = errors.New("invalid discipline id")

	// Shooting Event
	ErrInvalidShootingEventID = errors.New("invalid shooting event id")

	// Lane
	ErrLaneNameRequired = errors.New("lane name is required")
	ErrLaneAlreadyExists = errors.New("lane already exists")
	ErrAcademyBuildingNotFound = errors.New("academy building not found")
	ErrInvalidLaneID = errors.New("invalid lane id")
	ErrLaneAlreadyOccupied = errors.New("lane already occupied")

	// Session
	ErrPlayerAlreadyHasActiveSession = errors.New("player already has an active session")

	// Roles
	ErrRoleNotFound = errors.New("role not found")
	ErrInvalidRole = errors.New("not a valid role")

	// Users
	ErrUserNotFound = errors.New("error not found")

	// Profiles
	ErrProfileNotFound = errors.New("profile not found")


	// ----------------------------------------------------------

	// Validation Errors
	ErrInvalidEmailFormat = errors.New(
		"invalid email format",
	)

	ErrInvalidPhoneNumber = errors.New(
		"invalid phone number",
	)

	ErrWeakPassword = errors.New(
		"password must contain uppercase, lowercase, number and minimum 8 characters",
	)



	ErrInvalidUUID = errors.New(
		"invalid id",
	)


	ErrForbiddenDistrict = errors.New(
		"you are not authorized to manage this district admin",
	)

	// Business Errors
	// ErrEmailExists = errors.New(
	// 	"email already exists",
	// )

	

	ErrInvalidState = errors.New(
		"invalid assigned state",
	)

	ErrStateAdminNotFound = errors.New(
		"state admin not found",
	)

	ErrStateNotFound = errors.New(
		"state not found",
	)

	// Profile Related Errors
	ErrProfileAlreadyExists = errors.New(
		"profile with this email already exists",
	)

	ErrDistrictAdminNotFound = errors.New(
		"district admin not found",
	)


	ErrDistrictCoachNotFound = errors.New(
		"district coach Not found",
	)

	ErrInvalidDistrict = errors.New("invalid district")

	//Academy
	ErrInvalidAcademy = errors.New("invalid academy id")

	// Academy Admin
	ErrAcademyAdminNotFound = errors.New("academy admin not found")

)