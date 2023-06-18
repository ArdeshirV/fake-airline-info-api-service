package server

import (
	"fmt"
	"log"
	"github.com/labstack/echo/v4"
)

const (
	Normal      = "\033[0m"
	BoldBlue    = "\033[1;34m"
	BoldGreen   = "\033[1;32m"
	BoldYellow  = "\033[1;33m"
	BoldMagneta = "\033[1;35m"
)

var e *echo.Echo

func init() {
	e = echo.New()
}

func StartServer(portNumber string) {
	port := fmt.Sprintf(":%s", portNumber)
	loadTitle(port)
	rootRoutes(e)
	log.Fatal(e.Start(port))
}

func loadTitle(port string) {
	e.HideBanner = true
	defer fmt.Printf(Normal)

	title := "\n\t%sFake Airline Information Service!\n"
	fmt.Printf(title, BoldMagneta)

	message := "\t%shttp server started on: %shttp://localhost%s\n\n"
	fmt.Printf(message, BoldYellow, BoldBlue, port)
}
