package server

import (
	"fmt"
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/fake-airline-info-service/service"
)

const (
	ParamTime = "time"
	ParamCityA = "city_a"
	ParamCityB = "city_b"
	ParamCommand = "command"
	ParamFlightNo = "flightno"
	ParamLogoName = "logo_name"
	IndexFileName = "data/index.html"
)

func rootRoutes(e *echo.Echo) {
	e.GET("/", rootHandler)
	e.GET("/airline", findAirlineLogoHandler)
	e.GET("/flights", listFlightsHandler)
	e.GET("/airplanes", listAirplanesHandler)
	e.GET("/cities", listCitiesHandler)
	e.GET("/departure_dates", listDepartureDatesHandler)
	e.POST("/reserve_flight", listReserveFlightHandler)
}

func rootHandler(ctx echo.Context) error {
	return ctx.File(IndexFileName)
}

func findAirlineLogoHandler(ctx echo.Context) error {
	airlineLogoName := ctx.QueryParam(ParamLogoName)
	if airlineLogoName == "" {
		err := fmt.Errorf("the '%v' parameter is required", ParamLogoName)
		return echoErrorAsJSON(ctx, http.StatusBadRequest, err)
	}
	fileName, err := service.GetAirplaneLogoFileName(airlineLogoName)
	if err != nil {
		return echoErrorAsJSON(ctx, http.StatusBadRequest, err)
	}
	return ctx.File(fileName)
}

func listReserveFlightHandler(ctx echo.Context) error {
	command := ctx.QueryParam(ParamCommand)
	flightNo := ctx.QueryParam(ParamFlightNo)
	if command != "" || flightNo != "" {
		errMsg := ""
		dataIsNotEnough := false
		if command == "" {
			dataIsNotEnough = true
			errMsg += "\"command\" parameter is not defined correctly. "
		}
		if flightNo == "" {
			dataIsNotEnough = true
			errMsg += "\"flightno\" parameter is not defined correctly. "
		}
		if dataIsNotEnough {
			return echoStringAsJSON(ctx, http.StatusBadRequest, errMsg)
		} else {
			msg, err := service.SetRemainingCapacity(flightNo, command)
			if err != nil {
				return echoErrorAsJSON(ctx, http.StatusInternalServerError, err)
			}
			return echoJSON(ctx, http.StatusOK, msg)
		}
	}
	return echoStringAsJSON(ctx, http.StatusBadRequest, "bad request!")
}

func listAirplanesHandler(ctx echo.Context) error {
	data, err := service.GetAirplanes()
	if err != nil {
		return echoErrorAsJSON(ctx, http.StatusInternalServerError, err)
	}
	return echoJSON(ctx, http.StatusOK, data)
}

func listFlightsHandler(ctx echo.Context) error {
	flightNo := ctx.QueryParam(ParamFlightNo)
	if flightNo != "" {
		fliteredFlight, err := service.GetFlightsByFlightNo(flightNo)
		if err != nil {
			return echoErrorAsJSON(ctx, http.StatusInternalServerError, err)
		}
		return echoJSON(ctx, http.StatusOK, fliteredFlight)
	} else {
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
				return echoStringAsJSON(ctx, http.StatusBadRequest, errMsg)
			} else {
				filteredFlights, err := service.GetFlightsFromA2B(timeD, cityA, cityB)
				if err != nil {
					return echoErrorAsJSON(ctx, http.StatusBadRequest, err)
				}
				return echoJSON(ctx, http.StatusOK, filteredFlights)
			}
		} else {
			data, err := service.GetFlights()
			if err != nil {
				return echoErrorAsJSON(ctx, http.StatusInternalServerError, err)
			}
			return echoJSON(ctx, http.StatusOK, data)
		}
	}
	return echoStringAsJSON(ctx, http.StatusBadRequest, "bad request!")
}

func listCitiesHandler(ctx echo.Context) error {
	data, err := service.GetCities()
	if err != nil {
		return echoErrorAsJSON(ctx, http.StatusInternalServerError, err)
	}
	return echoJSON(ctx, http.StatusOK, data)
}

func listDepartureDatesHandler(ctx echo.Context) error {
	data, err := service.GetDepartureDates()
	if err != nil {
		return echoErrorAsJSON(ctx, http.StatusInternalServerError, err)
	}
	return echoJSON(ctx, http.StatusOK, data)
}
