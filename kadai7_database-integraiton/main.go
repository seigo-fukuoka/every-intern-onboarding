// きゅるりんってしてみてスケジュールアプリ - データベース統合版
// 機能: Webページから情報取得 + データベース保存 + API提供 + ブラウザ画面表示
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	migrate "github.com/rubenv/sql-migrate"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Event - イベント情報を格納する構造体
type Event struct {
	ID          int               `json:"id"`          // 自動生成ID
	Date        string            `json:"date"`        // 日付 (YYYY-MM-DD)
	Title       string            `json:"title"`       // イベント名
	IsAttending bool              `json:"isAttending"` // 参加予定フラグ
	Details     map[string]string `json:"details"`     // 追加情報（JSON形式）
}

// main - アプリケーションのエントリーポイント
func main() {
	// Webサーバー初期化（Echo = Go用のWebフレームワーク）
	e := echo.New()

	// データベース接続
	db, err := connectDB()
	if err != nil {
		panic(err)
	}
	// マイグレーション実行（テーブル作成・更新）
	err = initDB(db)
	if err != nil {
		panic(err)
	}

	// CORS設定（ブラウザからのアクセス許可設定）
	// React（ポート5174）からGo（ポート1323）へのアクセスを許可
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"}, // React開発サーバーのURL
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// ヘルスチェック用エンドポイント
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Echo!")
	})

	// GET /events - 全イベント取得API（ブラウザ画面用）
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

	// GET /scrape/all-events - Webページ情報取得実行API
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

	// Webサーバー起動（ポート1323）
	e.Logger.Fatal(e.Start(":1323"))
}

// scrapeAllEvent - きゅるりんってしてみてWebサイトからイベント情報をスクレイピング
// 重複チェックを行い、新しいイベントのみDBに保存
func scrapeAllEvent(db *sql.DB) ([]Event, error) {
	// Webページ読み取りツール初期化（iPhoneのフリをして情報取得）
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1"),
	)

	var newEvents []Event

	// Webページ解析: HTML文書からイベント情報を探し出す
	// 「article.tribe-events-calendar-month-mobile-events__mobile-event」というタグを見つけたら以下を実行
	// → きゅるりんサイトで1つのイベント情報が入っているHTMLの箱
	c.OnHTML("article.tribe-events-calendar-month-mobile-events__mobile-event", func(e *colly.HTMLElement) {

		// イベントの日時を取得（例: "2025.09.18 / 6:45 PM - 9:00 PM"）
		// HTMLの中から日時が書いてある部分を探して文字として取得
		dateTimeText := e.ChildText(".tribe-events-calendar-month-mobile-events__mobile-event-datetime")

		// イベントのタイトルを取得
		// HTMLの中からタイトルが書いてある部分を探して文字として取得
		titleText := e.ChildText(".tribe-events-calendar-month-mobile-events__mobile-event-title")

		// 見つけたイベント情報をコンソールに表示
		fmt.Printf("Found event: DateTime=%s, Title=%s\n", dateTimeText, titleText)

		// 日付部分だけを切り出し（"2025.09.18 / 6:45 PM" → "2025.09.18"）
		datePart := strings.Split(dateTimeText, " /")[0]
		datePart = strings.TrimSpace(datePart) // 前後の空白を削除

		// 日付の区切り文字を統一（"2025.09.18" → "2025-09-18"）
		formattedDate := strings.ReplaceAll(datePart, ".", "-")

		// 重複チェック: 同じ日付・同じタイトルのイベントが既にデータベースにあるか確認
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM events WHERE date = ? AND title = ?",
			formattedDate, titleText).Scan(&count)
		if err != nil {
			fmt.Printf("Error checking duplicate: %v\n", err)
			return
		}

		if count > 0 {
			// 既に同じイベントがあるので、このイベントは保存しない
			fmt.Printf("スキップ: %s (%s) - すでに存在\n", titleText, formattedDate)
			return
		}

		// イベントの詳細説明を取得
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

		// 新しいイベント情報をデータベースのeventsテーブルに1行追加
		result, err := db.Exec("INSERT INTO events (date, title, is_attending, details) VALUES (?, ?, ?, ?)",
			newEvent.Date, newEvent.Title, newEvent.IsAttending, string(detailsJSON))
		if err != nil {
			fmt.Printf("Error inserting event: %v\n", err)
			return
		}

		// データベースが自動で作ったIDを取得
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			fmt.Printf("Error getting last insert ID: %v\n", err)
			return
		}
		newEvent.ID = int(lastInsertID)

		newEvents = append(newEvents, newEvent)
		fmt.Printf("追加: %s (%s) - ID: %d\n", newEvent.Title, newEvent.Date, newEvent.ID)

	})

	// スクレイピングエラーハンドリング
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// 実際にWebページにアクセスして情報を取得開始
	err := c.Visit("https://www.kyurushite.com/schedule/")
	if err != nil {
		return nil, fmt.Errorf("failed to visit URL: %v", err)
	}

	return newEvents, nil
}

// connectDB - MySQLデータベースへの接続
func connectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/events_db?parseTime=true")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("データベース接続成功！")
	return db, nil
}

// initDB - データベース初期化（マイグレーション実行）
func initDB(db *sql.DB) error {
	return runMigrations(db)
}

// runMigrations - データベースの構造変更を実行
// migrations/フォルダ内の.sqlファイル（テーブル作成命令書）を順番に適用
func runMigrations(db *sql.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations/",
	}
	n, err := migrate.Exec(db, "mysql", migrations, migrate.Up)
	if err != nil {
		return err
	}
	fmt.Printf("Applied %d migrations!\n", n)
	return nil
}
