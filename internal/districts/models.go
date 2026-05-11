package districts

type District struct {
	ID 			 string `json:"id"`
	StateID 	 string `json:"state_id"`
	DistrictName string `json:"name"`
}

type DistrictResponse struct {
	ID   string    
	Name string 
}

type GetDistrictQueryParams struct {
	Search string `query:"search"`
	Order  string `query:"order"`
}