package pincode

type Pincode struct {
	ID int `json:"id"`

	Code string `json:"code"`

	DistrictID int `json:"district_id"`

	DistrictName string `json:"district_name"`

	StateName string `json:"state_name"`
}