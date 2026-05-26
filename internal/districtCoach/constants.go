package districtCoach

var AllowedDistrictCoachSortFields = map[string]string{
	"first_name":                 "p.first_name",
	"last_name":                  "p.last_name",
	"email":                      "p.email",
	"district_id":                "dc.district_id",
	"coach_code":                 "dc.coach_code",
	"created_at":                 "dc.created_at",
}