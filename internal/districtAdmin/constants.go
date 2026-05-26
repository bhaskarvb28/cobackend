package districtAdmin

var AllowedDistrictAdminSortFields = map[string]string{
	"first_name":  "p.first_name",
	"last_name":   "p.last_name",
	"email":       "p.email",
	"state_id":    "da.state_id",
	"district_id": "da.district_id",
	"created_at":  "da.created_at",
}