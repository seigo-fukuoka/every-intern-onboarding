package shared

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

func ScrapeAllEvent(db *sql.DB) ([]Event, error) {
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
