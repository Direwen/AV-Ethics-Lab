package session

type CreateSessionInput struct {
	AgeRange          int    `json:"age_range" validate:"required,min=1,max=6"`
	Gender            int    `json:"gender" validate:"required,min=1,max=4"`
	Country           string `json:"country" validate:"required"`
	Occupation        string `json:"occupation"`
	DrivingExperience int    `json:"driving_experience" validate:"min=1,max=3"`
	Fingerprint       string `json:"fingerprint" validate:"required"`
	SelfReportedNew   bool   `json:"self_reported_new"`
}

type SessionFeedback struct {
}
