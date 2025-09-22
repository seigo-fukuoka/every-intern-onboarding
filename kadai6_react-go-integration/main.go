package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Event struct {
	ID          int               `json:"id"`
	Date        string            `json:"date"`
	Title       string            `json:"title"`
	IsAttending bool              `json:"isAttending"`
	Details     map[string]string `json:"details"` // 柔軟な追加情報（全て文字列）
}

var (
	events = []Event{
		{
			ID: 1, Date: "2025-08-31", Title: "@JAM EXPO 2025 supported by UP-T", IsAttending: false,
			Details: map[string]string{
				"venue": "幕張メッセ",
				"time":  "開場16:00 / 開演17:00",
				"price": "前売り8,800円",
				"url":   "https://www.jam-expo.jp/",
			},
		},
		{
			ID: 2, Date: "2025-09-15", Title: "環やねpresents ＢＳ日テレぽっかる「ジャムムをきゅるして♡」ﾔﾈﾁｮｷお披露目会", IsAttending: false,
			Details: map[string]string{
				"venue":   "BS日テレスタジオ",
				"time":    "19:00-20:00",
				"channel": "BS日テレ",
				"note":    "生放送",
			},
		},
		// 一時的にコメントアウト（スクレイピングテスト用）
		// {
		// 	ID: 3, Date: "2025-09-18", Title: "ナカヨシファミリア〜きゅるりんってしてみて×ChumToto〜", IsAttending: false,
		// 	Details: map[string]string{
		// 		"venue": "渋谷CLUB QUATTRO",
		// 		"time":  "開場18:00 / 開演19:00",
		// 		"price": "前売り4,500円 / 当日5,000円",
		// 	},
		// },
		{
			ID: 4, Date: "2025-09-25", Title: "Kyururin♡Clinic 〜No.3ファントムシータ〜", IsAttending: false,
			Details: map[string]string{
				"venue":   "新宿ReNY",
				"time":    "開場18:30 / 開演19:30",
				"price":   "前売り5,500円",
				"presale": "ファンクラブ先行: 9/1-9/7",
			},
		},
		{
			ID: 5, Date: "2025-10-09", Title: "DEARSTAGE SHOWCASE 2025 AUTUMN", IsAttending: false,
			Details: map[string]string{
				"venue": "Zepp DiverCity",
				"time":  "開場17:00 / 開演18:00",
				"price": "前売り6,800円",
				"note":  "複数組出演",
			},
		},
		{
			ID: 6, Date: "2025-10-10", Title: "Unlock Secret Heart♡", IsAttending: false,
			Details: map[string]string{
				"venue": "恵比寿LIQUIDROOM",
				"time":  "開場18:00 / 開演19:00",
				"price": "前売り5,000円",
			},
		},
		{
			ID: 7, Date: "2025-10-20", Title: "あむはアイドル5ねんせい！〜みんなまとめてあむ心中〜", IsAttending: false,
			Details: map[string]string{
				"venue": "中野サンプラザ",
				"time":  "開場16:30 / 開演17:30",
				"price": "SS席8,000円 / S席6,500円",
				"note":  "5周年記念ライブ",
			},
		},
		{
			ID: 8, Date: "2025-11-08", Title: "Unlock Secret Heart♡", IsAttending: false,
			Details: map[string]string{
				"venue": "横浜Bay Hall",
				"time":  "開場17:30 / 開演18:30",
				"price": "前売り5,500円",
			},
		},
		{
			ID: 9, Date: "2026-01-17", Title: "GAKUSEI RUNWAY 2025 WINTER", IsAttending: false,
			Details: map[string]string{
				"venue":     "東京ビッグサイト",
				"time":      "13:00-17:00",
				"type":      "ファッションショー",
				"admission": "無料",
			},
		},
	}
	eventIDCounter = len(events)
)

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"}, // React開発サーバーのURL
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Echo!")
	})

	// GET /events エンドポイント: すべてのイベントを取得
	e.GET("/events", func(c echo.Context) error {
		return c.JSON(http.StatusOK, events) // 検索ロジックを削除し、すべてのイベントを返す
	})

	// 単一日付スクレイピングエンドポイント
	e.GET("/scrape/events", func(c echo.Context) error {
		dateParam := c.QueryParam("date")
		if dateParam == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Date parameter 'date' is required (YYYY-MM-DD)",
			})
		}

		// 日付フォーマットのバリデーション
		parsedDate, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Invalid date format. Use YYYY-MM-DD",
				"error":   err.Error(),
			})
		}

		scrapedEvent, err := scrapeEvent(parsedDate.Format("2006.01.02"))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Failed to scrape event",
				"error":   err.Error(),
			})
		}

		if scrapedEvent == nil {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": fmt.Sprintf("No event found for date %s", dateParam),
			})
		}
		return c.JSON(http.StatusOK, scrapedEvent)
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

		if newEvent.Title == "" || newEvent.Date == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Title and Date cannot be empty",
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
		if updateEvent.Title == "" || updateEvent.Date == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Title and Date cannot be empty",
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
		events[eventIndex].IsAttending = updateEvent.IsAttending
		events[eventIndex].Details = updateEvent.Details

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

// scrapeEvent は指定された日付のイベントをスクレイピングして返します
func scrapeEvent(targetDate string) (*Event, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"),
	)

	var foundEvent *Event
	// currentID := eventIDCounter + 1 // デバッグ中は一旦コメントアウト

	// 正しいセレクタでイベント要素を検索
	fmt.Printf("Setting up OnHTML handler for article selector\n")
	c.OnHTML("article.tribe-events-calendar-month-mobile-events__mobile-event", func(e *colly.HTMLElement) {
		fmt.Printf("OnHTML handler triggered for article element\n")
		if foundEvent != nil { // 既にイベントが見つかっている場合はスキップ
			return
		}

		// 日付情報を取得（例: "2025.09.18 / 6:45 PM - 9:00 PM"）
		dateTimeText := e.ChildText(".tribe-events-calendar-month-mobile-events__mobile-event-datetime")

		// タイトル情報を取得
		titleText := e.ChildText(".tribe-events-calendar-month-mobile-events__mobile-event-title")

		fmt.Printf("Found event: DateTime=%s, Title=%s\n", dateTimeText, titleText)

		// 日付部分だけを抽出（"2025.09.18"）
		datePart := strings.Split(dateTimeText, " /")[0]
		datePart = strings.TrimSpace(datePart)

		// 日付フォーマットを変換（"2025.09.18" → "2025-09-18"）
		formattedDate := strings.ReplaceAll(datePart, ".", "-")

		fmt.Printf("Extracted date: %s, Target date: %s\n", formattedDate, targetDate)

		// targetDateもYYYY-MM-DD形式に変換して比較
		normalizedTargetDate := strings.ReplaceAll(targetDate, ".", "-")
		fmt.Printf("Normalized target date: %s\n", normalizedTargetDate)

		if formattedDate == normalizedTargetDate {
			currentID := eventIDCounter + 1
			foundEvent = &Event{
				ID:          currentID,
				Date:        formattedDate,
				Title:       titleText,
				IsAttending: false,
				Details: map[string]string{
					"source":   "scraped",
					"url":      "https://www.kyurushite.com/schedule/",
					"datetime": dateTimeText, // 元の日時情報も保存
				},
			}
			eventIDCounter++
			fmt.Printf("Successfully created event: %+v\n", foundEvent)
		}
	})

	// エラーハンドリング
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// URLにアクセス
	err := c.Visit("https://www.kyurushite.com/schedule/")
	if err != nil {
		return nil, err
	}

	return foundEvent, nil
}
