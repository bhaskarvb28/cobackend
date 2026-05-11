package profiles

type CreateProfileInput struct {
	FirstName     string
	LastName      string
	Email         string
	PasswordHash  string
	ContactNumber string
	RoleID        string
}