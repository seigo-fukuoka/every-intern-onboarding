package main

import (
	"net/http"

	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Event struct {
	ID          int    `json:"id"`
	Date        string `json:"date"`
	Category    string `json:"category"`
	Title       string `json:"title"`
	IsAttending bool   `json:"isAttending"`
}

var (
	events = []Event{
		{ID: 1, Date: "2025-08-31", Category: "live", Title: "@JAM EXPO 2025 supported by UP-T", IsAttending: false},
		{ID: 2, Date: "2025-09-15", Category: "fan-meeting", Title: "環やねpresents ＢＳ日テレぽっかる「ジャムムをきゅるして♡」ﾔﾈﾁｮｷお披露目会", IsAttending: false},
		{ID: 3, Date: "2025-09-18", Category: "live", Title: "ナカヨシファミリア〜きゅるりんってしてみて×ChumToto〜", IsAttending: false},
		{ID: 4, Date: "2025-09-25", Category: "live", Title: "Kyururin♡Clinic 〜No.3ファントムシータ〜", IsAttending: false},
		{ID: 5, Date: "2025-10-09", Category: "live", Title: "DEARSTAGE SHOWCASE 2025 AUTUMN", IsAttending: false},
		{ID: 6, Date: "2025-10-10", Category: "live", Title: "Unlock Secret Heart♡", IsAttending: false},
		{ID: 7, Date: "2025-10-20", Category: "live", Title: "あむはアイドル5ねんせい！〜みんなまとめてあむ心中〜", IsAttending: false},
		{ID: 8, Date: "2025-11-08", Category: "live", Title: "Unlock Secret Heart♡", IsAttending: false},
		{ID: 9, Date: "2026-01-17", Category: "others", Title: "GAKUSEI RUNWAY 2025 WINTER", IsAttending: false},
	}
	eventIDCounter = len(events)
)

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"}, // React開発サーバーのURL
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Echo!")
	})

	e.GET("/events", func(c echo.Context) error {
		searchQuery := c.QueryParam("q")

		if searchQuery == "" {
			return c.JSON(http.StatusOK, events)
		}

		filteredEvents := []Event{}
		for _, event := range events {
			if strings.Contains(strings.ToLower(event.Title), strings.ToLower(searchQuery)) ||
				strings.Contains(strings.ToLower(event.Category), strings.ToLower(searchQuery)) {
				filteredEvents = append(filteredEvents, event)
			}
		}

		return c.JSON(http.StatusOK, filteredEvents)
	})

	// POST /events エンドポイント: 新しいイベントを追加
	e.POST("/events", func(c echo.Context) error {
		newEvent := new(Event)

		if err := c.Bind(newEvent); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Invalid request body",
				"error":   err.Error(),
			})
		}

		if newEvent.Title == "" || newEvent.Date == "" || newEvent.Category == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Title, Date, and Category cannot be empty",
			})
		}

		eventIDCounter++
		newEvent.ID = eventIDCounter

		events = append(events, *newEvent)

		return c.JSON(http.StatusCreated, newEvent)
	})

	// GET /events/:id エンドポイント: 特定のイベントを取得
	e.GET("/events/:id", func(c echo.Context) error {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Invalid event ID",
				"error":   err.Error(),
			})
		}

		for _, event := range events {
			if event.ID == id {
				return c.JSON(http.StatusOK, event)
			}
		}

		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "Event not found",
		})
	})

	e.PUT("/events/:id", func(c echo.Context) error {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Invalid event ID",
				"error":   err.Error(),
			})
		}

		updateEvent := new(Event)
		if err := c.Bind(updateEvent); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Invalid request body",
				"error":   err.Error(),
			})
		}
		if updateEvent.Title == "" || updateEvent.Date == "" || updateEvent.Category == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Title, Date, and Category cannot be empty",
			})
		}

		eventIndex := -1
		for i, event := range events {
			if event.ID == id {
				eventIndex = i
				break
			}
		}

		if eventIndex == -1 {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "Event not found",
			})
		}

		events[eventIndex].Title = updateEvent.Title
		events[eventIndex].Date = updateEvent.Date
		events[eventIndex].Category = updateEvent.Category
		events[eventIndex].IsAttending = updateEvent.IsAttending

		return c.JSON(http.StatusOK, events[eventIndex])
	})

	e.DELETE("/events/:id", func(c echo.Context) error {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Invalid event ID",
				"error":   err.Error(),
			})
		}

		eventIndex := -1
		for i, event := range events {
			if event.ID == id {
				eventIndex = i
				break
			}
		}
		if eventIndex == -1 {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "Event not found",
			})
		}

		events = append(events[:eventIndex], events[eventIndex+1:]...)
		return c.NoContent(http.StatusNoContent)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
