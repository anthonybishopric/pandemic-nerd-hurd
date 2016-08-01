package pandemic

type DiseaseType string

type DiseaseData struct {
	Type          DiseaseType `json:"type"`
	Incurable     bool        `json:"incurable,omitempty"`
	Untreatable   bool        `json:"untreatable,omitempty"`
	BecomingFaded bool        `json:"becoming_faded,omitempty"`
}

var Yellow = DiseaseData{
	Type: DiseaseType("Yellow"),
}
var Blue = DiseaseData{
	Type:          DiseaseType("Blue"),
	Incurable:     true, // TODO: make configurable with a gamestate
	Untreatable:   true,
	BecomingFaded: true,
}
var Red = DiseaseData{
	Type: DiseaseType("Red"),
}
var Black = DiseaseData{
	Type: DiseaseType("Black"),
}
var Faded = DiseaseData{
	Type:          DiseaseType("Faded"),
	Incurable:     true,
	Untreatable:   true,
	BecomingFaded: true,
}
