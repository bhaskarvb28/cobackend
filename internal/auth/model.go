package auth

type RegisterInput struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	Password      string `json:"password"` // plain password
	RoleID        int    `json:"role_id"`
	ContactNumber string `json:"contact_number"`
}