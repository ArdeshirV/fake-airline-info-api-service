package server

import (
	"fmt"
	"net/http"
	"github.com/labstack/echo/v4"
	models "github.com/the-go-dragons/fake-airline-info-service/domain"
	"github.com/the-go-dragons/fake-airline-info-service/service"
	"github.com/the-go-dragons/fake-airline-info-service/config"
)

const (
	ParamTime = "time"
	ParamCityA = "city_a"
	ParamCityB = "city_b"
	ParamCommand = "command"
	CommandReturn = "return"
	CommandReserve = "reserve"
	ParamFlightNo = "flightno"
	RootPage =
`<html>
	<header>
		<title>Fake Airline Information Service API</title>
	</header>
	<body style="color:white; background-color:black;" text-align:left; " >
		<head>
			<h1 style="color:yellow; text-align:center;">Fake Airline Information Service API</h1>
		<head>
		<section>
			<pre>
				<a href="http://localhost:3000" target="_blank"  rel="noopener noreferrer">http://localhost:3000</a>
				<a href="http://localhost:3000/cities" target="_blank" rel="noopener noreferrer">http://localhost:3000/cities</a>
				<a href="http://localhost:3000/airplanes" target="_blank" rel="noopener noreferrer">http://localhost:3000/airplanes</a>
				<a href="http://localhost:3000/departure_dates" target="_blank" rel="noopener noreferrer">http://localhost:3000/departure_dates</a>
				<a href="http://localhost:3000/flights" target="_blank" rel="noopener noreferrer">http://localhost:3000/flights</a>
				<a href="http://localhost:3000/flights?flightno=EAX254" target="_blank" rel="noopener noreferrer">http://localhost:3000/flights?flightno=EAX254</a>
				<a href="http://localhost:3000/flights?city_a=New%20York&city_b=Paris&time=2023-06-14" target="_blank" rel="noopener noreferrer">http://localhost:3000/flights?city_a=New%20York&city_b=Paris&time=2023-06-14</a>
			</pre>
			<pre style="color:green;">
				curl -X GET "http://localhost:3000/cities"
				curl -X GET "http://localhost:3000/airplanes"
				curl -X GET "http://localhost:3000/departure_dates"
				curl -X GET "http://localhost:3000/flights"
				curl -X GET "http://localhost:3000/flights?flightno=EAX254"
				curl -X GET "http://localhost:3000/flights?city_a=New%20York&city_b=Paris&time=2023-06-14"
				<br/>
				curl -X POST "http://localhost:3000/reserve_flight?flightno=EAX254&command=return"
				curl -X POST "http://localhost:3000/reserve_flight?flightno=EAX254&command=reserve"
			</pre>
		</section>
		<footer style="text-align:center;">
			<p><a target="_blank" rel="noopener noreferrer" href="https://github.com/the-go-dragons/fake-airline-info-api-service">Visit the Project in Github</a></p>
			<p style="text-align: center; width: 100%; ">Copyright&copy; 2023 <a href="https://github.com/the-go-dragons">The Go Dragons Team</a>, Licensed under MIT</p>
		</footer>
	<body>
</html>`
)

func rootRoutes(e *echo.Echo) {
	e.GET("/", rootHandler)
	e.GET("/flights", listFlightsHandler)
	e.GET("/airplanes", listAirplanesHandler)
	e.GET("/cities", listCitiesHandler)
	e.GET("/departure_dates", listDepartureDatesHandler)
	e.POST("/reserve_flight", listReserveFlightHandler)
}

func rootHandler(ctx echo.Context) error {
	return ctx.HTML(http.StatusOK, RootPage)
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
			return ctx.HTML(http.StatusBadRequest, htmlErrorMsgString(errMsg))
		} else {
			msg, err := service.SetRemainingCapacity(flightNo, command)
			return echoJSON(ctx, http.StatusOK, msg)
		}
	}
	return ctx.HTML(http.StatusBadRequest, htmlErrorMsgString("Bad request!"))
}

func listAirplanesHandler(ctx echo.Context) error {
	data, err := service.GetAirplanes()
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, htmlErrorMsg(err))
	}
	return echoJSON(ctx, http.StatusOK, data)
}

func listFlightsHandler(ctx echo.Context) error {
	flightNo := ctx.QueryParam(ParamFlightNo)
	if flightNo != "" {
		fliteredFlight, err := service.GetFlightsByFlightNo(flightNo)
		if err != nil {
			return ctx.HTML(http.StatusInternalServerError, htmlErrorMsg(err))
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
				return ctx.HTML(http.StatusBadRequest, htmlErrorMsgString(errMsg))
			} else {
				filteredFlights, err := service.GetFlightsFromA2B(timeD, cityA, cityB)
				if err != nil {
					return ctx.HTML(http.StatusBadRequest, htmlErrorMsg(err))
				}
				return echoJSON(ctx, http.StatusOK, filteredFlights)
			}
		} else {
			data, err := service.GetFlights()
			if err != nil {
				return ctx.HTML(http.StatusInternalServerError, htmlErrorMsg(err))
			}
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

func htmlErrorMsg(err error) string {
	return fmt.Sprintf("<strong style=\"color:red\">Error: %v</strong>", err)
}

func htmlErrorMsgString(errMsg string) string {
	return fmt.Sprintf("<strong style=\"color:red\">Error: %v</strong>", errMsg)
}

func echoJSON(ctx echo.Context, status int, data interface{}) error {
	if config.IsDebugModeEnabled() {
		return ctx.JSONPretty(status, data, "    ")
	} else {
		return ctx.JSON(status, data)
	}
}
