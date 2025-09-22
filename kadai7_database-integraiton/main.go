package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
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

func main() {
	e := echo.New()

	db, err := connectDB()
	if err != nil {
		panic(err)
	}
	err = initDB(db)
	if err != nil {
		panic(err)
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:5175"}, // React開発サーバーのURL
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Echo!")
	})

	// GET /events エンドポイント: すべてのイベントを取得
	e.GET("/events", func(c echo.Context) error {
		rows, err := db.Query("SELECT id, date, title, is_attending, details FROM events ORDER BY date ASC")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Failed to fetch events",
				"error":   err.Error(),
			})
		}
		defer rows.Close()

		var events []Event
		for rows.Next() {
			var event Event
			var detailsJSON string
			err := rows.Scan(&event.ID, &event.Date, &event.Title, &event.IsAttending, &detailsJSON)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"message": "Failed to scan event",
					"error":   err.Error(),
				})
			}

			// JSON文字列をmap[string]stringに変換
			err = json.Unmarshal([]byte(detailsJSON), &event.Details)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"message": "Failed to parse event details",
					"error":   err.Error(),
				})
			}

			events = append(events, event)
		}

		return c.JSON(http.StatusOK, events)
	})

	// GET /scrape/all-events エンドポイント: すべてのイベントをスクレイピングして追加
	e.GET("/scrape/all-events", func(c echo.Context) error {
		newEvents, err := scrapeAllEvent(db)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Failed to scrape all events",
				"error":   err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": fmt.Sprintf("%d件の新しいイベントを追加しました", len(newEvents)),
			"count":   len(newEvents),
			"events":  newEvents,
		})

	})

	e.Logger.Fatal(e.Start(":1323"))
}

// scrapeAllEvent は指定された日付のイベントをスクレイピングして返します
func scrapeAllEvent(db *sql.DB) ([]Event, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1"),
	)

	var newEvents []Event
	// currentID := eventIDCounter + 1 // デバッグ中は一旦コメントアウト

	// 正しいセレクタでイベント要素を検索
	fmt.Printf("Setting up OnHTML handler for article selector\n")
	c.OnHTML("article.tribe-events-calendar-month-mobile-events__mobile-event", func(e *colly.HTMLElement) {
		fmt.Printf("OnHTML handler triggered for article element\n")

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

		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM events WHERE date = ? AND title = ?",
			formattedDate, titleText).Scan(&count)
		if err != nil {
			fmt.Println("Error checking duplicate: %v\n", err)
			return
		}

		if count > 0 {
			fmt.Printf("スキップ: %s (%s) - すでに存在\n", titleText, formattedDate)
			return
		}

		fullDescription := e.ChildText("div.tribe-events-single-event-description")

		newEvent := Event{
			Date:        formattedDate,
			Title:       titleText,
			IsAttending: false,
			Details:     map[string]string{"description": fullDescription},
		}

		detailsJSON, err := json.Marshal(newEvent.Details)
		if err != nil {
			fmt.Printf("Error marshalling details: %v\n", err)
			return
		}

		// DBにINSERT（IDはAUTO_INCREMENTで自動生成）
		result, err := db.Exec("INSERT INTO events (date, title, is_attending, details) VALUES (?, ?, ?, ?)",
			newEvent.Date, newEvent.Title, newEvent.IsAttending, string(detailsJSON))
		if err != nil {
			fmt.Printf("Error inserting event: %v\n", err)
			return
		}

		// 生成されたIDを取得してnewEventに設定
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			fmt.Printf("Error getting last insert ID: %v\n", err)
			return
		}
		newEvent.ID = int(lastInsertID)

		newEvents = append(newEvents, newEvent)
		fmt.Printf("追加: %s (%s) - ID: %d\n", newEvent.Title, newEvent.Date, newEvent.ID)

	})

	// エラーハンドリング
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err := c.Visit("https://www.kyurushite.com/schedule/")
	if err != nil {
		return nil, fmt.Errorf("failed to visit URL: %v", err)
	}

	return newEvents, nil
}

func connectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/events_db")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("データベース接続成功！")
	return db, nil
}

func initDB(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS events (
		id INT AUTO_INCREMENT PRIMARY KEY,
		date VARCHAR(10) NOT NULL,
		title VARCHAR(255) NOT NULL,
		is_attending BOOLEAN DEFAULT FALSE,
		details JSON,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		KEY idx_date (date)
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}
	fmt.Println("eventsテーブル作成完了！")
	return nil
}

func migrateMemoryToDB(db *sql.DB, events []Event) error {
	for _, event := range events {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM events WHERE date = ? AND title = ?", event.Date, event.Title).Scan(&count)
		if err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		detailsJSON, err := json.Marshal(event.Details)
		if err != nil {
			return err
		}
		_, err = db.Exec("INSERT INTO events (date, title, is_attending, details) VALUES (?, ?, ?, ?)", event.Date, event.Title, event.IsAttending, string(detailsJSON))
		if err != nil {
			return err
		}
	}

	return nil
}
