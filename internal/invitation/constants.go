package invitation

var InvitationPermissions = map[string][]string{

	"super_admin": {
		"state_admin",
	},

	"state_admin": {
		"district_admin",
		"district_coach",
	},

	"district_admin": {
		"academy_admin",
	},

	"academy_admin": {
		"academy_coach",
		"player",
	},
}