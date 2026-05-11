package stateadmin

type InviteStateAdminInput struct {
	Email           string `json:"email"`
	AssignedStateID string `json:"assigned_state_id"`
}

// type CreateStateAdminInput struct {
// 	FirstName     string `json:"first_name"`
// 	LastName      string `json:"last_name"`
// 	Email         string `json:"email"`
// 	Password      string `json:"password"`
// 	ContactNumber string `json:"contact_number"`
// 	AssignedState string    `json:"assigned_state"`
// }



// type UpdateAssignedStateInput struct {
// 	AssignedState string `json:"assigned_state"`
// }

// type StateAdminResponse struct {
// 	ID              string `json:"id"`
// 	FirstName       string `json:"first_name"`
// 	LastName        string `json:"last_name"`
// 	Email           string `json:"email"`
// 	ContactNumber   string `json:"contact_number"`
// 	AssignedState   int    `json:"assigned_state"`
// }

// type GetStateAdminsQuery struct {
// 	Page           int
// 	Limit          int
// 	Search         string
// 	AssignedState  string
// }