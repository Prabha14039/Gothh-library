package main

import (
	"Gothh/templates"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Count struct {
	Count int
}

func main(){

	darkMode := true

	var items []int =  []int {1,2,3,4,5,6}
	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/", func (c echo.Context) error{
		return templates.Index(darkMode,items).Render(c.Request().Context(),c.Response().Writer)
	})

	e.GET("/home", func (c echo.Context) error{
		return templates.Home(darkMode,items).Render(c.Request().Context(),c.Response().Writer)
	})

	e.Static("/css","css")
	e.Static("/static","static")
	e.Logger.Fatal(e.Start(":8080"))

}
