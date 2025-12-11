package dto

type CreateSessionInput struct {
	Demographics    map[string]any `json:"demographics"`
	Fingerprint     string         `json:"fingerprint" validate:"required"`
	SelfReportedNew bool           `json:"self_reported_new"`
}
