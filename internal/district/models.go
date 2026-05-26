package district

// District represents a district entity.
type District struct {

	// ID is the unique identifier of the district.
	ID int `json:"id"`

	// StateID is the identifier of the
	// state the district belongs to.
	StateID int `json:"state_id"`

	// DistrictName is the name of the district.
	DistrictName string `json:"name"`
}

// DistrictResponse represents the public
// response structure for district data.
type DistrictResponse struct {

	// ID is the unique identifier of the district.
	ID int `json:"id"`

	// Name is the display name of the district.
	Name string `json:"name"`
}

// GetDistrictQueryParams contains query parameters
// used for fetching districts.
type GetDistrictQueryParams struct {

	// Search filters districts by name.
	Search string `query:"search"`

	// Order defines the sorting order
	// of the returned districts.
	Order string `query:"order"`
}