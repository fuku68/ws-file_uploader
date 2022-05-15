package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"echo-upload/handler"
)

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	api := e.Group("/api")
	h := handler.NewHandler()
	h.Register(api)
	e.Logger.Fatal(e.Start(":1323"))
}
