package shared

import "errors"

var (

	// Generic Errors
	ErrInvalidRequestBody = errors.New(
		"invalid request body",
	)

	ErrUnauthorized = errors.New(
		"unauthorized",
	)

	ErrForbidden = errors.New(
		"forbidden",
	)

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

	ErrInvalidCredentials = errors.New("invalid email or password")


	ErrInvalidUUID = errors.New(
		"invalid id",
	)

	// Business Errors
	ErrEmailExists = errors.New(
		"email already exists",
	)

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
)

