package academy

// AllowedAcademySortFields maps API sort query parameters (keys)
// to their corresponding database columns (values) to prevent SQL injection.
// The "a." prefix references the "academies" table alias, and the "d." prefix
// references the "districts" table alias in joined queries.
var AllowedAcademySortFields = map[string]string{
	"id":          "a.id",
	"name":        "a.name",
	"state_id":    "d.state_id",
	"district_id": "a.district_id",
	"created_at":  "a.created_at",
	"updated_at":  "a.updated_at",
}