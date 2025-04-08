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
	e := echo.New()
	e.Use(middleware.Logger())

	count := Count {Count :0}

	e.GET("/", func (c echo.Context) error{
		return templates.Index(count.Count).Render(c.Request().Context(),c.Response().Writer)
	})

	e.GET("/count", func (c echo.Context) error{
		count.Count++
		return templates.Counter(count.Count).Render(c.Request().Context(),c.Response().Writer)
	})

	e.Static("/css","css")
	e.Logger.Fatal(e.Start(":8080"))

}
