package llm

type ScenarioRequest struct {
	TemplateName   string
	GridDimensions string
	GridData       []string
	Factors        map[string]interface{}
}

type ScenarioResponse struct {
	Narrative string
	Entities  []RawEntity
	Factors   map[string]interface{}
}

type RawEntity struct {
	Type string
	Row  int
	Col  int
	Meta RawEntityMeta
}

type RawEntityMeta struct {
	Behavior string
	Occluded bool
}
