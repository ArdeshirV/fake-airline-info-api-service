package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	models "github.com/the-go-dragons/fake-airline-info-service/domain"
	"github.com/the-go-dragons/fake-airline-info-service/config"
)

type Timestamp time.Time

const (
	Normal  = "\033[0m"
	Red     = "\033[0;31m"
	BoldRed = "\033[1;31m"
)

const (
	TimeLayout = "2006-01-02"
)

// TODO: Incompleted implementation
func SetRemainingCapacity(flightNo, command string) (string, error) {
	msg := ""
	flights, err := GetFlights()
	if err != nil {
		return msg, errorHandler("Failed to read flight data", err)
	}
	filteredFlights := make([]models.Flight, 0)
	for _, flight := range flights {
		if flight.FlightNo == flightNo {
			filteredFlights = append(filteredFlights, flight)
		}
	}
	if len(filteredFlights) > 0 {
		selectedFlight := filteredFlights[0]
		msg = fmt.Sprintf("%v", selectedFlight)
	}
	return fmt.Sprintf(`[{"message": "%v"}]`, msg), nil
}

func GetAirplanes() ([]models.Airplane, error) {
	flights, err := GetFlights()
	if err != nil {
		return nil, errorHandler("Failed to read flight data", err)
	}
	airplanes := make([]models.Airplane, len(flights))
	for _, flight := range flights {
		airplanes = append(airplanes, flight.Airplane)
	}
	return airplanes, nil
}

func GetFlights() ([]models.Flight, error) {
	flightData, err := os.ReadFile("data/flight.json")
	if err != nil {
		return nil, errorHandler("Failed to read flight data", err)
	}
	var flights []models.Flight
	err = json.Unmarshal(flightData, &flights)
	if err != nil {
		return nil, errorHandler("Failed to unmarshal flight data", err)
	}
	return flights, nil
}

func GetFlightsByFlightNo(flightNo string) ([]models.Flight, error) {
	flights, err := GetFlights()
	if err != nil {
		return nil, errorHandler("Failed to get flight data", err)
	}
	filteredFlights := make([]models.Flight, 0)
	for _, flight := range flights {
		if flight.FlightNo == flightNo {
			filteredFlights = append(filteredFlights, flight)
		}
	}
	return filteredFlights, nil
}

func GetFlightsFromA2B(timeD, cityA, cityB string) ([]models.Flight, error) {
	flights, err := GetFlights()
	if err != nil {
		return nil, errorHandler("Failed to get flight data", err)
	}
	filteredFlights := make([]models.Flight, 0)
	specifiedDate, err := time.Parse(TimeLayout, timeD)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to parse \"time\". %v", err)
		return nil, errorHandler(errMsg, err)
	}
	for _, flight := range flights {
		if AreDatesEqual(flight.DepartureTime, specifiedDate) &&
		Normalize(cityA) == Normalize(flight.Departure.City.Name) &&
		Normalize(cityB) == Normalize(flight.Destination.City.Name) {
			filteredFlights = append(filteredFlights, flight)
		}
	}
	return filteredFlights, nil
}

func GetCities() ([]models.City, error) {
	appendUniqueCity := func(slice []models.City, data models.City) []models.City {
		for _, value := range slice {
			if value == data {
				return slice
			}
		}
		return append(slice, data)
	}
	flights, err := GetFlights()
	if err != nil {
		return nil, errorHandler("Failed to read city data", err)
	}
	cities := make([]models.City, 0)
	for _, flight := range flights {
		cities = appendUniqueCity(cities, flight.Departure.City)
		cities = appendUniqueCity(cities, flight.Destination.City)
	}
	return cities, nil
}

func GetDepartureDates() ([]time.Time, error) {
	flights, err := GetFlights()
	if err != nil {
		return nil, errorHandler("Failed to read departure-time data", err)
	}
	departureTimes := make([]time.Time, len(flights))
	for _, flight := range flights {
		departureTimes = append(departureTimes, flight.DepartureTime)
	}
	return departureTimes, nil
}

func errorHandler(message string, err error) error {
	errMsgHTML := fmt.Sprintf("%s %v\n", message, err)
	if config.IsDebugModeEnabled() {
		errMsgColor := fmt.Sprintf("Error: %s%s\n%s%v%s\n", BoldRed, message, Red, err, Normal)
		log.New(os.Stderr, "\n", 1).Print(errMsgColor)
	}
	return errors.New(errMsgHTML)
}

func (t *Timestamp) UnmarshalParam(src string) error {
	ts, err := time.Parse(time.RFC3339, src)
	*t = Timestamp(ts)
	return err
}

func Normalize(text string) string {
	return strings.ToLower(strings.TrimSpace(text))
}

func AreDatesEqual(a, b time.Time) bool {
	return a.Day() == b.Day() && a.Month() == b.Month() && a.Year() == b.Year()
}
