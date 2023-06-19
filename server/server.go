package server

import (
	"fmt"
	"log"
	"github.com/labstack/echo/v4"
	colors "github.com/the-go-dragons/fake-airline-info-service/config/colors"
)

var e *echo.Echo

func init() {
	e = echo.New()
}

func StartServer(portNumber string) {
	loadTitle(portNumber)
	rootRoutes(e)
	log.Fatal(e.Start(fmt.Sprintf(":%s", portNumber)))
}

func loadTitle(port string) {
	e.HideBanner = true
	defer fmt.Print(colors.Normal)
	title := "\n\t%sFake Airline Information Service!\n"
	fmt.Printf(title, colors.BoldMagneta)
	message := "\t%shttp server started on: %shttp://localhost:%s\n\n"
	fmt.Printf(message, colors.BoldYellow, colors.BoldBlue, port)
}
