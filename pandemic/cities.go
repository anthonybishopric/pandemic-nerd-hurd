package pandemic

type PanicLevel int

func (p PanicLevel) CanBuildResearchStations() bool {
	return int(p) < 2
}

const (
	Nothing = PanicLevel(iota)
	Unstable
	Rioting2
	Rioting3
	Collapsing
	Fallen
)

type DiseaseType struct {
	Color         string `json:"color"`
	Incurable     bool   `json:"incurable,omitempty"`
	Untreatable   bool   `json:"untreatable,omitempty"`
	BecomingFaded bool   `json:"becoming_faded,omitempty"`
	HasZombies    bool
}

var Yellow = DiseaseType{
	Color: "Yellow",
}
var Blue = DiseaseType{
	Color:         "Blue",
	Incurable:     true, // TODO: make configurable with a gamestate
	Untreatable:   true,
	BecomingFaded: true,
}
var Red = DiseaseType{
	Color: "Red",
}
var Black = DiseaseType{
	Color: "Black",
}
var Faded = DiseaseType{
	Color:       "Faded",
	HasZombies:  true,
	Incurable:   true,
	Untreatable: true,
}

type City struct {
	Name        string
	Epidemic    bool
	FundedEvent bool
	Disease     DiseaseType
	PanicLevel  PanicLevel
}

var SanFrancisco = City{
	Name:    "San Francisco",
	Disease: Blue,
}
var Washington = City{
	Name:    "Washington",
	Disease: Blue,
}
var Atlanta = City{
	Name:    "Atlanta",
	Disease: Blue,
}
var Montreal = City{
	Name:    "Montreal",
	Disease: Blue,
}
var Chicago = City{
	Name:    "Chicago",
	Disease: Blue,
}
var NewYork = City{
	Name:    "New York",
	Disease: Blue,
}
var London = City{
	Name:    "London",
	Disease: Blue,
}
var Essen = City{
	Name:    "Essen",
	Disease: Blue,
}
var StPetersburg = City{
	Name:    "St. Petersburg",
	Disease: Blue,
}
var Milan = City{
	Name:    "Milan",
	Disease: Blue,
}
var Paris = City{
	Name:    "Paris",
	Disease: Blue,
}
var Madrid = City{
	Name:    "Madrid",
	Disease: Blue,
}

var LosAngeles = City{
	Name:    "Los Angeles",
	Disease: Yellow,
}

var Miami = City{
	Name:    "Miami",
	Disease: Yellow,
}
var MexicoCity = City{
	Name:    "Mexico City",
	Disease: Yellow,
}
var Bogota = City{
	Name:    "Bogota",
	Disease: Yellow,
}
var Lima = City{
	Name:    "Lima",
	Disease: Yellow,
}
var Santiago = City{
	Name:    "Santiago",
	Disease: Yellow,
}
var SaoPaulo = City{
	Name:    "Sao Paulo",
	Disease: Yellow,
}
var BuenosAires = City{
	Name:    "Buenos Aires",
	Disease: Yellow,
}
var Lagos = City{
	Name:    "Lagos",
	Disease: Yellow,
}
var Khartoum = City{
	Name:    "Khartoum",
	Disease: Yellow,
}
var Kinshasa = City{
	Name:    "Kinshasa",
	Disease: Yellow,
}
var Johannesburg = City{
	Name:    "Johannesburg",
	Disease: Yellow,
}

var Algiers = City{
	Name:    "Algiers",
	Disease: Black,
}
var Istanbul = City{
	Name:    "Istanbul",
	Disease: Black,
}
var Cairo = City{
	Name:    "Cairo",
	Disease: Black,
}
var Riydah = City{
	Name:    "Riydah",
	Disease: Black,
}
var Baghdad = City{
	Name:    "Baghdad",
	Disease: Black,
}
var Moscow = City{
	Name:    "Moscow",
	Disease: Black,
}
var Tehran = City{
	Name:    "Tehran",
	Disease: Black,
}
var Delhi = City{
	Name:    "Delhi",
	Disease: Black,
}
var Karachi = City{
	Name:    "Karachi",
	Disease: Black,
}
var Mumbai = City{
	Name:    "Mumbai",
	Disease: Black,
}
var Kolkata = City{
	Name:    "Kolkata",
	Disease: Black,
}
var Chennai = City{
	Name:    "Chennai",
	Disease: Black,
}

var Beijing = City{
	Name:    "Beijing",
	Disease: Red,
}
var Seoul = City{
	Name:    "Seoul",
	Disease: Red,
}
var Tokyo = City{
	Name:    "Tokyo",
	Disease: Red,
}
var Shanghai = City{
	Name:    "Shanghai",
	Disease: Red,
}
var Taipei = City{
	Name:    "Taipei",
	Disease: Red,
}
var Osaka = City{
	Name:    "Osaka",
	Disease: Red,
}
var HongKong = City{
	Name:    "HongKong",
	Disease: Red,
}
var Bangkok = City{
	Name:    "Bangkok",
	Disease: Red,
}
var HoChiMinhCity = City{
	Name:    "HoChiMinhCity",
	Disease: Red,
}
var Jakarta = City{
	Name:    "Jakarta",
	Disease: Red,
}
var Manila = City{
	Name:    "Manila",
	Disease: Red,
}
var Sydney = City{
	Name:    "Sydney",
	Disease: Red,
}

var AllCities = []City{
	SanFrancisco,
	Washington,
	Atlanta,
	Montreal,
	Chicago,
	NewYork,
	London,
	Essen,
	StPetersburg,
	Milan,
	Paris,
	Madrid,
	LosAngeles,
	Miami,
	MexicoCity,
	Bogota,
	Lima,
	Santiago,
	SaoPaulo,
	BuenosAires,
	Lagos,
	Khartoum,
	Kinshasa,
	Johannesburg,
	Algiers,
	Istanbul,
	Cairo,
	Riydah,
	Baghdad,
	Moscow,
	Tehran,
	Delhi,
	Karachi,
	Mumbai,
	Kolkata,
	Chennai,
	Beijing,
	Seoul,
	Tokyo,
	Shanghai,
	Taipei,
	Osaka,
	HongKong,
	Bangkok,
	HoChiMinhCity,
	Jakarta,
	Manila,
	Sydney,
}

func AllCitiesWithDisease(disease DiseaseType) []City {
	cities := []City{}
	for _, city := range AllCities {
		if city.Disease == disease {
			cities = append(cities, city)
		}
	}
	return cities
}
