package districts

type District struct {
	ID 			 int 	`json:"id"`
	StateID 	 int 	`json:"state_id"`
	DistrictName string `json:"district_name"`
}

type DistrictResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}