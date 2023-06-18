package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	models "github.com/the-go-dragons/fake-airline-info-service/domain"
)

const (
	Normal  = "\033[0m"
	Red     = "\033[0;31m"
	BoldRed = "\033[1;31m"
)

func GetAirplanes() ([]models.Airplane, error) {
	flights, err := GetFlights()
	if err != nil {
		return nil, errorHandler("Failed to read <Airplane(s)> data", err)
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
		fmt.Println(err)
		return nil, errorHandler("Failed to unmarshal flight data", err)
	}

	return flights, nil
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
	errMsgColor := fmt.Sprintf("Error: %s%s\n%s%v%s\n", BoldRed, message, Red, err, Normal)
	log.New(os.Stderr, "\n", 1).Print(errMsgColor)

	return errors.New(errMsgHTML)
}
