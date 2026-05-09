package stateadmin

type CreateStateAdminInput struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	ContactNumber string `json:"contact_number"`
	AssignedState int    `json:"assigned_state"`
}

type UpdateAssignedStateInput struct {
	AssignedState int `json:"assigned_state"`
}

type StateAdminResponse struct {
	ID              string `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	ContactNumber   string `json:"contact_number"`
	AssignedState   int    `json:"assigned_state"`
}

type GetStateAdminsQuery struct {
	Page           int
	Limit          int
	Search         string
	AssignedState  int
}