package pandemic

type DiseaseType string

type DiseaseData struct {
	Type             DiseaseType `json:"type"`
	Incurable        bool        `json:"incurable,omitempty"`
	Untreatable      bool        `json:"untreatable,omitempty"`
	BecomingFaded    bool        `json:"becoming_faded,omitempty"`
	InfectOnCityDraw bool        `json:"infect_on_city_draw,omitempty"`
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
	Type:             DiseaseType("Faded"),
	Incurable:        true,
	Untreatable:      true,
	BecomingFaded:    true,
	InfectOnCityDraw: true,
}

func (dt DiseaseType) String() string {
	return string(dt)
}

var diseaseDataMap map[DiseaseType]DiseaseData

func init() {
	diseaseDataMap = map[DiseaseType]DiseaseData{
		Yellow.Type: Yellow,
		Blue.Type:   Blue,
		Red.Type:    Red,
		Black.Type:  Black,
		Faded.Type:  Faded,
	}
}

func DataForDisease(dt DiseaseType) DiseaseData {
	return diseaseDataMap[dt]
}

func CurableDiseases() []DiseaseType {
	ret := []DiseaseType{}
	for dt, data := range diseaseDataMap {
		if !data.Incurable {
			ret = append(ret, dt)
		}
	}
	return ret
}
