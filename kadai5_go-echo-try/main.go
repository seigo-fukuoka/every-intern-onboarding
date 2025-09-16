package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Task struct {
	Name string `json:"name"`
	Status string `json:"status"`
}

var tasks = []Task{
	{Name: "ご飯食べる", Status: "完了"},
	{Name: "家に帰る", Status: "未着手"},
	{Name: "寝る", Status: "未着手"},
}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Echo!")
	})

	e.GET("/tasks", func(c echo.Context) error {
		return c.JSON(http.StatusOK, tasks)
	})

	e.Logger.Fatal(e.Start(":1323"))
}