package academyAdmin

var AllowedAcademyAdminSortFields = map[string]string{
	"first_name":  "p.first_name",
	"last_name":   "p.last_name",
	"email":       "p.email",
	"state_id":    "aa.state_id",
	"district_id": "aa.district_id",
	"academy_id":  "aa.academy_id",
	"created_at":  "aa.created_at",
}