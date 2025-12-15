package model

// Age range codes
const (
	AgeRange18to24 = 1
	AgeRange25to34 = 2
	AgeRange35to44 = 3
	AgeRange45to54 = 4
	AgeRange55to64 = 5
	AgeRange65Plus = 6
)

// Gender codes
const (
	GenderMale        = 1
	GenderFemale      = 2
	GenderNonBinary   = 3
	GenderPreferNotTo = 4
)

// Driving experience codes
const (
	DrivingLicensed = 1
	DrivingNone     = 2
	DrivingLearner  = 3
)

type Session struct {
	BaseModel
	// For the Unique Anonymity Protocol
	Fingerprint     string `gorm:"type:varchar(255)" json:"fingerprint"`
	IsDuplicate     bool   `gorm:"type:boolean;default:true" json:"is_duplicate"`
	SelfReportedNew bool   `gorm:"type:boolean;default:false" json:"self_reported_new"`
	// Demographics
	AgeRange          int    `gorm:"type:smallint;not null" json:"age_range"`
	Gender            int    `gorm:"type:smallint;not null" json:"gender"`
	Country           string `gorm:"type:varchar(10);not null" json:"country"` // ISO code
	Occupation        string `gorm:"type:varchar(50)" json:"occupation"`
	DrivingExperience int    `gorm:"type:smallint" json:"driving_experience"`
}
