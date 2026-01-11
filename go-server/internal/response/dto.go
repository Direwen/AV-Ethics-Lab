package response

type SubmitResponseInput struct {
	RankingOrder   []string `json:"ranking_order" validate:"required,len=3"`
	ResponseTimeMs int64    `json:"response_time_ms" validate:"required,min=0"`
	IsTimeout      bool     `json:"is_timeout"`
	HasInteracted  bool     `json:"has_interacted"`
}

type SubmitResponseOutput struct {
	Response   *Response `json:"response"`
	IsComplete bool      `json:"is_complete"`
}

type RankResult struct {
	Outcome string `json:"outcome"`
	Count   int64  `json:"count"`
}

type FactorCount struct {
	FactorValue string `json:"factor_value"`
	Outcome     string `json:"outcome"`
	Count       int64  `json:"count"`
}

type ResponseTimeMSCount struct {
	Second int64 `json:"second"`
	Count  int64 `json:"count"`
}
