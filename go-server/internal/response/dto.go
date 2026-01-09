package response

type SubmitResponseInput struct {
	RankingOrder   []string `json:"ranking_order" validate:"required,len=3"`
	ResponseTimeMs int64    `json:"response_time_ms" validate:"required,min=0"`
	IsTimeout      bool     `json:"is_timeout"`
	HasInteracted  bool     `json:"has_interacted"`
}
