package session

import "github.com/direwen/go-server/internal/shared/models"

// Re-export from shared models for backward compatibility
type Session = models.Session
type SessionStatus = models.SessionStatus

const (
	StatusActive    = models.StatusActive
	StatusCompleted = models.StatusCompleted
	StatusExpired   = models.StatusExpired
	StatusAbandoned = models.StatusAbandoned
)

var SessionStatusErrorMsg = models.SessionStatusErrorMsg

// Demographic constants
const (
	AgeRange18to24 = models.AgeRange18to24
	AgeRange25to34 = models.AgeRange25to34
	AgeRange35to44 = models.AgeRange35to44
	AgeRange45to54 = models.AgeRange45to54
	AgeRange55to64 = models.AgeRange55to64
	AgeRange65Plus = models.AgeRange65Plus
)

const (
	GenderMale        = models.GenderMale
	GenderFemale      = models.GenderFemale
	GenderNonBinary   = models.GenderNonBinary
	GenderPreferNotTo = models.GenderPreferNotTo
)

const (
	DrivingLicensed = models.DrivingLicensed
	DrivingNone     = models.DrivingNone
	DrivingLearner  = models.DrivingLearner
)
