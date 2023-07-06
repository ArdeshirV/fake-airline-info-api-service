package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/the-go-dragons/fake-airline-info-service/config"
	colors "github.com/the-go-dragons/fake-airline-info-service/config/colors"
	"github.com/the-go-dragons/fake-airline-info-service/config/logos"
	models "github.com/the-go-dragons/fake-airline-info-service/domain"
)

const (
	TimeLayout     = "2006-01-02"
	JSONFile       = "data/flight.json"
	CommandReturn  = "return"
	CommandReserve = "reserve"
)

type ReserveResponse struct {
	Message           string `json:"message"`
	FlightNo          string `json:"flightno"`
	Capacity          int    `json:"capacity"`
	RemainingCapacity int    `json:"remainingcapacity"`
}

func SetRemainingCapacity(flightNo string, cmd string) (string, error) {
	command := string(cmd)
	flights, err := GetFlights()
	if err != nil {
		return "", errorHandler("Failed to read flight data", err)
	}
	count := 0
	index := -1
	for i, flight := range flights {
		if Normalize(flight.FlightNo) == Normalize(flightNo) {
			count++
			index = i
		}
	}
	if count > 0 {
		if count > 1 {
			errMsgFmt := "Duplication detected. More than one Flight found FlightNo:%v"
			errMsg := fmt.Sprintf(errMsgFmt, flightNo)
			return "", errorHandler(errMsg, err)
		}
		msg := ""
		dataChanged := false
		cmd := Normalize(command)
		switch cmd {
		case CommandReserve:
			if flights[index].RemainingCapacity > 0 {
				flights[index].RemainingCapacity--
				dataChanged = true
				msg = "Remaning capacity reduced"
			} else {
				msg = "There is not any remaning capacity."
			}
		case CommandReturn:
			if int(flights[index].RemainingCapacity) < int(flights[index].Airplane.Capacity) {
				flights[index].RemainingCapacity++
				dataChanged = true
				msg = "Remaning capacity increased."
			} else {
				msg = "All capacity of airplane is available. You can not reduce remaining capacity anymore."
			}
		default:
			errMsg := fmt.Sprintf("The \"%v\" command is unkown.", command)
			return "", errorHandler(errMsg, err)
		}
		if dataChanged {
			err := setFlights(flights)
			if err != nil {
				return "", errorHandler("Failed to set flights data", err)
			}
			msg += " Remaining capacity updated with new value."
		}
		reserveResponse := ReserveResponse{
			Message:           msg,
			FlightNo:          flightNo,
			Capacity:          int(flights[index].Airplane.Capacity),
			RemainingCapacity: flights[index].RemainingCapacity,
		}
		bytes, err := json.Marshal(reserveResponse)
		if err != nil {
			return msg, errorHandler("Failed to convert ReserveResponse struct to JSON", err)
		}
		return string(bytes), nil
	}
	errMsg := "Failed to find specified Flight with FlightNo:%v to reserve/return"
	return "", errorHandler(fmt.Sprintf(errMsg, flightNo), err)
}

func GetAirplaneLogoFileName(name string) (string, error) {
	airlineLogo, err := logos.GetAirlineLogoByName(logos.AirlineName(name))
	if err != nil {
		return "", err
	}
	return string(airlineLogo), nil
}

func setFlights(flights []models.Flight) error {
	bytes, err := json.Marshal(flights)
	if err != nil {
		return errorHandler("Failed to convert []models.Flight to JSON", err)
	}
	err = os.WriteFile(JSONFile, bytes, 0664)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to write []models.Flight as JSON in file \"%s\"", JSONFile)
		return errorHandler(errMsg, err)
	}
	return nil
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
	flightData, err := os.ReadFile(JSONFile)
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
	errMsg := fmt.Sprintf("%s %v\n", message, err)
	if config.IsDebugMode() {
		errMsgColor := fmt.Sprintf("Error: %s%s\n%s%v%s\n",
			colors.BoldRed, message, colors.Red, err, colors.Normal)
		log.New(os.Stderr, "\n", 1).Print(errMsgColor)
	}
	return errors.New(errMsg)
}

func Normalize(text string) string {
	return strings.ToLower(strings.TrimSpace(text))
}

func AreDatesEqual(a, b time.Time) bool {
	return a.Day() == b.Day() && a.Month() == b.Month() && a.Year() == b.Year()
}
