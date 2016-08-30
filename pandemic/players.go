package pandemic

type CharacterType int

const (
	Medic = CharacterType(iota)
	Dispatcher
	Researcher
	Scientist
	Civilian
	QuarantineSpecialist
	Colonel
	OperationsExpert
	Generalist
)

type Player struct {
	HumanName string `json:"human_name"`
	Character *Character
	Location  CityName
	Cities    []CityCard
}

type Character struct {
	Name string        `json:"name"`
	Type CharacterType `json:"character_type"`
}
