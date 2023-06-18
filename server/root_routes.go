package server

import (
	"fmt"
	"net/http"
	"time"
	"strings"

	"github.com/labstack/echo/v4"
	models "github.com/the-go-dragons/fake-airline-info-service/domain"
	"github.com/the-go-dragons/fake-airline-info-service/service"
	"github.com/the-go-dragons/fake-airline-info-service/config"
)

type Timestamp time.Time

const (
	ParamTime = "time"
	ParamCityA = "city_a"
	ParamCityB = "city_b"
	ParamReturn = "return"
	ParamReserve = "reserve"
	ParamFlightNo = "flightno"
	TimeLayout = "2006-01-02"
	// TODO: This string should be replaced with README.md file and it's contents should be raw HTML not markdown
	RootPage = "<strong>Fake Airline Information Service API</strong><br/><br/><a target=\"_blank\" rel=\"noopener noreferrer\" href=\"https://github.com/the-go-dragons/fake-airline-info-api-service\">Check the project in Github</a>"
)

func rootRoutes(e *echo.Echo) {
	e.GET("/", rootHandler)
	e.GET("/flights", listFlightsHandler)
	e.GET("/airplanes", listAirplanesHandler)
	e.GET("/cities", listCitiesHandler)
	e.GET("/departure_dates", listDepartureDatesHandler)
	e.POST("/reserve", listReserveHandler)
}

func rootHandler(ctx echo.Context) error {
	return ctx.HTML(http.StatusOK, RootPage)
}

func listReserveHandler(ctx echo.Context) error {
	flights, err := service.GetFlights()
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, htmlErrorMsg(err))
	}
	flightNo := ctx.QueryParam(ParamFlightNo)
	if flightNo != "" {
		filteredFlights := make([]models.Flight, 0)
		for _, flight := range flights {
			if flight.FlightNo == flightNo {
				filteredFlights = append(filteredFlights, flight)
			}
		}
		if len(filteredFlights) > 0 {
			selectedFlight := filteredFlights[0]
			// The solution of Ticket-4 goes here
			//
			return echoJSON(ctx, http.StatusOK, selectedFlight)
		}
	}

	return echoJSON(ctx, http.StatusOK, "Nothing")
}

func listAirplanesHandler(ctx echo.Context) error {
	data, err := service.GetAirplanes()
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, htmlErrorMsg(err))
	}

	return echoJSON(ctx, http.StatusOK, data)
}

func listFlightsHandler(ctx echo.Context) error {
	data, err := service.GetFlights()
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, htmlErrorMsg(err))
	}

	flightNo := ctx.QueryParam(ParamFlightNo)
	if flightNo != "" {
		filteredFlights := make([]models.Flight, 0)
		for _, flight := range data {
			if flight.FlightNo == flightNo {
				filteredFlights = append(filteredFlights, flight)
			}
		}
		return echoJSON(ctx, http.StatusOK, filteredFlights)
	} else {
		// The solution of Ticket-2 goes here
		cityA := ctx.QueryParam(ParamCityA)
		cityB := ctx.QueryParam(ParamCityB)
		timeD := ctx.QueryParam(ParamTime)
		if timeD != "" || cityA != "" || cityB != "" {
			errMsg := ""
			dataIsNotEnough := false
			if timeD == "" {
				dataIsNotEnough = true
				errMsg += "\"time\" is not defined correctly. "
			}
			if cityA == "" {
				dataIsNotEnough = true
				errMsg += "\"city_a\" is not defined correctly. "
			}
			if cityB == "" {
				dataIsNotEnough = true
				errMsg += "\"city_b\" is not defined correctly. "
			}

			if dataIsNotEnough {
				return ctx.HTML(http.StatusInternalServerError, htmlErrorMsgString(errMsg))
			} else {
				// Search to find all flights from cityA to cityB in specified date-time
				filteredFlights := make([]models.Flight, 0)
				specifiedDate, err := time.Parse(TimeLayout, timeD)
				if err != nil {
					errMsg := fmt.Sprintf("Failed to parse \"time\". %v", err)
					return ctx.HTML(http.StatusInternalServerError, htmlErrorMsgString(errMsg))
				}
				for _, flight := range data {
					if AreDatesEqual(flight.DepartureTime, specifiedDate) &&
					normalize(cityA) == normalize(flight.Departure.City.Name) &&
					normalize(cityB) == normalize(flight.Destination.City.Name) {
						filteredFlights = append(filteredFlights, flight)
					}
				}
				if len(filteredFlights) > 0 {
					return echoJSON(ctx, http.StatusOK, filteredFlights)
				}
				// Not found
				return echoJSON(ctx, http.StatusOK, []models.Flight{})
			}
		} else {
			// Ticket-1 Solution
			return echoJSON(ctx, http.StatusOK, data)
		}
	}

	return ctx.HTML(http.StatusBadRequest, htmlErrorMsgString("Bad request!"))
}

func listCitiesHandler(ctx echo.Context) error {
	data, err := service.GetCities()
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, htmlErrorMsg(err))
	}

	return echoJSON(ctx, http.StatusOK, data)
}

func listDepartureDatesHandler(ctx echo.Context) error {
	data, err := service.GetDepartureDates()
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, htmlErrorMsg(err))
	}

	return echoJSON(ctx, http.StatusOK, data)
}

func (t *Timestamp) UnmarshalParam(src string) error {
	ts, err := time.Parse(time.RFC3339, src)
	*t = Timestamp(ts)
	return err
}

func htmlErrorMsg(err error) string {
	return fmt.Sprintf("<strong style=\"color:red\">Error: %v</strong>", err)
}

func htmlErrorMsgString(errMsg string) string {
	return fmt.Sprintf("<strong style=\"color:red\">Error: %v</strong>", errMsg)
}

func normalize(text string) string {
	return strings.ToLower(strings.Trim(text, " \t"))
}

func AreDatesEqual(a, b time.Time) bool {
	return a.Day() == b.Day() && a.Month() == b.Month() && a.Year() == b.Year()
}

func echoJSON(ctx echo.Context, status int, data interface{}) error {
	if config.IsDebugModeEnabled() {
		return ctx.JSONPretty(status, data, "    ")
	} else {
		return ctx.JSON(status, data)
	}
}
