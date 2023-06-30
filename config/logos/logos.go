package logos

import (
	"fmt"
	"strings"
)


// Public Interface ----------------------------------------------------------------------

type AirlineID int
type AirlineName string
type AirlineLogo string

const (
	AmericanAirline AirlineID = iota
	DeltaAirLine
	UnitedAirline
	Lufthansa
	Emirate
	BritishAirway
	AirFrance
	CathayPacificAirway
	QantasAirway
	SingaporeAirline
)

const (
	AmericanAirlineName AirlineName = "american_airlines"
	DeltaAirLineName AirlineName = "delta_air_lines"
	UnitedAirlineName AirlineName = "united_airlines"
	LufthansaName AirlineName = "lufthansa"
	EmirateName AirlineName = "emirates"
	BritishAirwayName AirlineName = "british_airways"
	AirFranceName AirlineName = "air_france"
	CathayPacificAirwayName AirlineName = "cathay_pacific_airways"
	QantasAirwayName AirlineName = "qantas_airways"
	SingaporeAirlineName AirlineName = "singapore_airlines"
)

func GetAirlineName(id AirlineID) AirlineName {
	load()
	return airlineNames[id]
}

func GetAirlineLogo(id AirlineID) AirlineLogo {
	load()
	return airlineLogos[id]
}

func GetAirlineNameList() []AirlineName {
	load()
	return airlineNames[:]
}

func GetAirlineLogoList() []AirlineLogo {
	load()
	return airlineLogos[:]
}

func GetAirlineLogoByName(name AirlineName) (AirlineLogo, error) {
	load()
	airName := normalize(name)
	for index, airlineName := range airlineNames {
		if airlineName == airName {
			return airlineLogos[index], nil
		}
	}
	return AirlineLogo(""), fmt.Errorf("airline name:'%v' not found", name)
}

func normalize(text AirlineName) AirlineName {
	return AirlineName(strings.ToLower(strings.TrimSpace(string(text))))
}

// Private Implementation ----------------------------------------------------------------

const (
	americanAirlinesLogo AirlineLogo = "img/airlines_logo/American_Airlines.jpeg"
	deltaAirLinesLogo AirlineLogo = "img/airlines_logo/Delta_Air_Lines.png"
	unitedAirlinesLogo AirlineLogo = "img/airlines_logo/United_Airlines.jpeg"
	lufthansaLogo AirlineLogo = "img/airlines_logo/Lufthansa.png"
	emiratesLogo AirlineLogo = "img/airlines_logo/Emirates.png"
	britishAirwaysLogo AirlineLogo = "img/airlines_logo/British_Airways.png"
	airFranceLogo AirlineLogo = "img/airlines_logo/Air_France.png"
	cathayPacificAirwaysLogo AirlineLogo = "img/airlines_logo/Cathay_Pacific_Airways.png"
	qantasAirwaysLogo AirlineLogo = "img/airlines_logo/Qantas_Airways.jpeg"
	singaporeAirlinesLogo AirlineLogo = "img/airlines_logo/Singapore_Airlines.jpg"
)

var (
	airlineNames []AirlineName
	airlineLogos []AirlineLogo
)

func load() {
	if airlineNames == nil {
		airlineNames = make([]AirlineName, 10)
		airlineNames[AmericanAirline] = AmericanAirlineName
		airlineNames[DeltaAirLine] = DeltaAirLineName
		airlineNames[UnitedAirline] = UnitedAirlineName
		airlineNames[Lufthansa] = LufthansaName
		airlineNames[Emirate] = EmirateName
		airlineNames[BritishAirway] = BritishAirwayName
		airlineNames[AirFrance] = AirFranceName
		airlineNames[CathayPacificAirway] = CathayPacificAirwayName
		airlineNames[QantasAirway] = QantasAirwayName
		airlineNames[SingaporeAirline] = SingaporeAirlineName
	}
	if airlineLogos == nil {
		airlineLogos = make([]AirlineLogo, 10)
		airlineLogos[AmericanAirline] = americanAirlinesLogo
		airlineLogos[DeltaAirLine] = deltaAirLinesLogo
		airlineLogos[UnitedAirline] = unitedAirlinesLogo
		airlineLogos[Lufthansa] = lufthansaLogo
		airlineLogos[Emirate] = emiratesLogo
		airlineLogos[BritishAirway] = britishAirwaysLogo
		airlineLogos[AirFrance] = airFranceLogo
		airlineLogos[CathayPacificAirway] = cathayPacificAirwaysLogo
		airlineLogos[QantasAirway] = qantasAirwaysLogo
		airlineLogos[SingaporeAirline] = singaporeAirlinesLogo
	}
}
