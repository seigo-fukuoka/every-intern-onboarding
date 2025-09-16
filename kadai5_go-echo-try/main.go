package main

import (
	"net/http"

	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type Task struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

var (
	tasks = []Task{
		{ID: 1, Name: "ご飯食べる", Status: "完了"},
		{ID: 2, Name: "家に帰る", Status: "未着手"},
		{ID: 3, Name: "寝る", Status: "未着手"},
	}
	taskIDCounter = len(tasks)
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Echo!")
	})

	e.GET("/tasks", func(c echo.Context) error {
		searchQuery := c.QueryParam("q")

		if searchQuery == "" {
			return c.JSON(http.StatusOK, tasks)
		}

		filteredTasks := []Task{}
		for _, task := range tasks {
			if strings.Contains(strings.ToLower(task.Name), strings.ToLower(searchQuery)) ||
				strings.Contains(strings.ToLower(task.Status), strings.ToLower(searchQuery)) {
				filteredTasks = append(filteredTasks, task)
			}
		}

		return c.JSON(http.StatusOK, filteredTasks)
	})

	// POST /tasks エンドポイント: 新しいタスクを追加
	e.POST("/tasks", func(c echo.Context) error {
		newTask := new(Task)

		if err := c.Bind(newTask); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Invalid request body",
				"error":   err.Error(),
			})
		}

		if newTask.Name == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Task name cannot be empty",
			})
		}

		taskIDCounter++
		newTask.ID = taskIDCounter

		tasks = append(tasks, *newTask)

		return c.JSON(http.StatusCreated, newTask)
	})

	// GET /tasks/:id エンドポイント: 特定のタスクを取得
	e.GET("/tasks/:id", func(c echo.Context) error {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Invalid task ID",
				"error":   err.Error(),
			})
		}

		for _, task := range tasks {
			if task.ID == id {
				return c.JSON(http.StatusOK, task)
			}
		}

		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "Task not found",
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
