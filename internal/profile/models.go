package profile

import (
	"cobackend/internal/auth"
	"time"
)


type CreateProfileInput struct {
	FirstName     string
	LastName      string
	Email         string
	PasswordHash  string
	ContactNumber string
	RoleID        string
}


type ProfileResponse struct {
	User    auth.UserResponse `json:"user"`
	Profile interface{}       `json:"profile"`
}

type StateAdminProfileResponse struct {
	ProfileCompleted bool   `json:"profile_completed"`
	StateID          int16  `json:"state_id"`
	StateName        string `json:"state_name"`
	DPDPConsent      bool   `json:"dpdp_consent"`
}

type DistrictAdminProfileResponse struct {
	ProfileCompleted bool   `json:"profile_completed"`
	DPDPConsent      bool   `json:"dpdp_consent"`

	DistrictID   int32  `json:"district_id"`
	DistrictName string `json:"district_name"`

	StateID   int16  `json:"state_id"`
	StateName string `json:"state_name"`
}

type Discipline struct {
	ID          int16  `json:"id"`
	Code        string `json:"code"`
	DisplayName string `json:"display_name"`
	IsPrimary   bool   `json:"is_primary,omitempty"`
}

type PincodeInfo struct {
	ID       int32  `json:"id"`
	Code     string `json:"code"`
	District string `json:"district"`
	State    string `json:"state"`
}

type AcademySummary struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Address string `json:"address"`

	District string `json:"district"`

	State string `json:"state"`
}

type DistrictCoachProfileResponse struct {
	ProfileCompleted bool   `json:"profile_completed"`
	DPDPConsent      bool   `json:"dpdp_consent"`

	DistrictID   int32  `json:"district_id"`
	DistrictName string `json:"district_name"`

	StateID   int16  `json:"state_id"`
	StateName string `json:"state_name"`

	CoachCode *string `json:"coach_code,omitempty"`
	CoachingCertificateProof *string `json:"coaching_certificate_proof,omitempty"`

	Disciplines []Discipline `json:"disciplines"`
}

type AcademyAdminProfileResponse struct {
	ProfileCompleted bool `json:"profile_completed"`
	DPDPConsent bool `json:"dpdp_consent"`

	GSTIN *string `json:"gstin,omitempty"`

	RegistrationProof *string `json:"registration_proof,omitempty"`

	AcademyID string `json:"academy_id"`

	AcademyName string `json:"academy_name"`

	AcademyAddress string `json:"academy_address"`

	DistrictID int32 `json:"district_id"`

	DistrictName string `json:"district_name"`

	StateID int16 `json:"state_id"`

	StateName string `json:"state_name"`
}

type AcademyCoachProfileResponse struct {
	ProfileCompleted bool `json:"profile_completed"`

	DPDPConsent bool `json:"dpdp_consent"`

	CoachCode *string `json:"coach_code,omitempty"`

	CoachingCertificateProof *string `json:"coaching_certificate_proof,omitempty"`

	AcademyID string `json:"academy_id"`

	AcademyName string `json:"academy_name"`

	AcademyAddress string `json:"academy_address"`

	DistrictID int32 `json:"district_id"`

	DistrictName string `json:"district_name"`

	StateID int16 `json:"state_id"`

	StateName string `json:"state_name"`

	Disciplines []Discipline `json:"disciplines"`
}

type PlayerGuardian struct {
	ID                 int64   `json:"id"`
	FullName           string  `json:"full_name"`
	Relationship       *string `json:"relationship,omitempty"`
	ContactNumber      string  `json:"contact_number"`
	AlternativeContact *string `json:"alternative_contact,omitempty"`
	ParentalConsent    bool    `json:"parental_consent"`
	IsPrimary          bool    `json:"is_primary"`
}

type PlayerPersonalInfo struct {
	DateOfBirth time.Time `json:"date_of_birth"`

	Gender string `json:"gender"`

	Nationality string `json:"nationality"`

	PlaceOfBirth *string `json:"place_of_birth,omitempty"`

	City *string `json:"city,omitempty"`

	ResidentialAddress *string `json:"residential_address,omitempty"`

	Pincode *PincodeInfo `json:"pincode,omitempty"`	
	
	Education *string `json:"education,omitempty"`

	InstitutionName *string `json:"institution_name,omitempty"`

	Occupation *string `json:"occupation,omitempty"`

	TemporarySportID *string `json:"temporary_sport_id,omitempty"`
}

type PlayerSportsProfile struct {

	UnitOfRepresentation *string `json:"unit_of_representation,omitempty"`

	DominantHand *string `json:"dominant_hand,omitempty"`

	HeightCM *float64 `json:"height_cm,omitempty"`

	WeightKG *float64 `json:"weight_kg,omitempty"`

	ShoeSize *string `json:"shoe_size,omitempty"`

	TracksuitSize *string `json:"tracksuit_size,omitempty"`
}

type PlayerPassport struct {
	PassportNumber *string `json:"passport_number,omitempty"`

	PassportIssueDate *time.Time `json:"passport_issue_date,omitempty"`

	PassportExpiryDate *time.Time `json:"passport_expiry_date,omitempty"`

	PassportIssuingAuthority *string `json:"passport_issuing_authority,omitempty"`

	PassportPlaceOfIssue *string `json:"passport_place_of_issue,omitempty"`
}

type PlayerProfileResponse struct {
	ProfileCompleted bool `json:"profile_completed"`

	DPDPConsent bool `json:"dpdp_consent"`

	Status string `json:"status"`

	JoinedAt time.Time `json:"joined_at"`

	Academy AcademySummary `json:"academy"`

	CurrentCoachUserID *string `json:"current_coach_user_id,omitempty"`

	PersonalInfo *PlayerPersonalInfo `json:"personal_info,omitempty"`

	SportsProfile *PlayerSportsProfile `json:"sports_profile,omitempty"`

	Disciplines []Discipline `json:"disciplines"`

	Passport *PlayerPassport `json:"passport,omitempty"`

	Guardians []PlayerGuardian `json:"guardians"`
}


type CompleteStateAdminProfileInput struct {
	DPDPConsent bool `json:"dpdp_consent"`
}

type CompleteDistrictAdminProfileInput struct {
	DPDPConsent bool `json:"dpdp_consent"`
}

type CompleteDistrictCoachProfileInput struct {
	DPDPConsent bool `json:"dpdp_consent"`

	CoachCode string `json:"coach_code"`

	CoachingCertificateProof string `json:"coaching_certificate_proof"`

	DisciplineIDs []int16 `json:"discipline_ids"`
}

type CompleteAcademyAdminProfileInput struct {
	DPDPConsent bool `json:"dpdp_consent"`

	GSTIN string `json:"gstin"`

	RegistrationProof string `json:"registration_proof"`
}

type CompleteAcademyCoachProfileInput struct {
	DPDPConsent bool `json:"dpdp_consent"`

	CoachCode string `json:"coach_code"`

	CoachingCertificateProof string `json:"coaching_certificate_proof"`

	DisciplineIDs []int16 `json:"discipline_ids"`
}

type PlayerPersonalInfoInput struct {
	DateOfBirth time.Time `json:"date_of_birth"`

	Gender string `json:"gender"`

	Nationality string `json:"nationality"`

	PlaceOfBirth *string `json:"place_of_birth,omitempty"`

	City *string `json:"city,omitempty"`

	ResidentialAddress *string `json:"residential_address,omitempty"`

	PincodeID *int32 `json:"pincode_id,omitempty"`

	Education *string `json:"education,omitempty"`

	InstitutionName *string `json:"institution_name,omitempty"`

	Occupation *string `json:"occupation,omitempty"`

	TemporarySportID *string `json:"temporary_sport_id,omitempty"`
}

type PlayerSportsProfileInput struct {
	UnitOfRepresentation *string `json:"unit_of_representation,omitempty"`

	DominantHand *string `json:"dominant_hand,omitempty"`

	HeightCM *float64 `json:"height_cm,omitempty"`

	WeightKG *float64 `json:"weight_kg,omitempty"`

	ShoeSize *string `json:"shoe_size,omitempty"`

	TracksuitSize *string `json:"tracksuit_size,omitempty"`
}

type PlayerDisciplineInput struct {
	DisciplineID int16 `json:"discipline_id"`

	IsPrimary bool `json:"is_primary"`
}

type PlayerGuardianInput struct {
	FullName string `json:"full_name"`

	Relationship *string `json:"relationship,omitempty"`

	ContactNumber string `json:"contact_number"`

	AlternativeContact *string `json:"alternative_contact,omitempty"`

	ParentalConsent bool `json:"parental_consent"`

	IsPrimary bool `json:"is_primary"`
}

type CompletePlayerProfileInput struct {
	DPDPConsent bool `json:"dpdp_consent"`

	PersonalInfo PlayerPersonalInfoInput `json:"personal_info"`

	SportsProfile PlayerSportsProfileInput `json:"sports_profile"`

	Disciplines []PlayerDisciplineInput `json:"disciplines"`

	Passport *PlayerPassport `json:"passport,omitempty"`

	Guardians []PlayerGuardianInput `json:"guardians"`
}