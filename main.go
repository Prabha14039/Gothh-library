package main

import (
	"Gothh/helpers"
	"Gothh/templates"
	"database/sql"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	dsn := helpers.FetchEnv()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	defer db.Close()

	darkMode := true

	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		images, err := helpers.Images_fetch(db)

		if err != nil {
			log.Println("Error fetching images:", err)
			return c.String(500, "Error fetching images")
		}

		return templates.Index(darkMode, images).Render(c.Request().Context(), c.Response().Writer)
	})

	e.GET("/home", func(c echo.Context) error {
		images, err := helpers.Images_fetch(db)

		if err != nil {
			log.Println("Error fetching images:", err)
			return c.String(500, "Error fetching images")
		}
		return templates.Home(darkMode, images).Render(c.Request().Context(), c.Response().Writer)
	})

	e.Static("/css", "css")
	e.Static("/static", "static")
	e.Logger.Fatal(e.Start(":8080"))

}
