package districts

type District struct {
	ID 			 int `json:"id"`
	StateID 	 int `json:"state_id"`
	DistrictName string `json:"name"`
}

type DistrictResponse struct {
	ID   int    
	Name string 
}

type GetDistrictQueryParams struct {
	Search string `query:"search"`
	Order  string `query:"order"`
}