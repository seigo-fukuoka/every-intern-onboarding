package service

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"

	"kadai7_database-integraiton/repository"
)

// EventService - イベントのドメイン層（ビジネスロジック）
type EventService struct {
	eventRepo *repository.EventRepository
}

// NewEventService - EventServiceのコンストラクタ
func NewEventService(eventRepo *repository.EventRepository) *EventService {
	return &EventService{
		eventRepo: eventRepo,
	}
}

// GetAllEvents - 全イベント取得（シンプルな委譲）
func (s *EventService) GetAllEvents() ([]repository.Event, error) {
	return s.eventRepo.GetAll()
}

// ScrapeAndSaveEvents - スクレイピング実行してDB保存（shared/scraper.goから移行）
func (s *EventService) ScrapeAndSaveEvents(limit int) ([]repository.Event, error) {
	// Webページ読み取りツール初期化（iPhoneのフリをして情報取得）
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1"),
	)

	var newEvents []repository.Event
	var stopped bool // limit到達フラグ

	// Webページ解析: HTML文書からイベント情報を探し出す
	c.OnHTML("article.tribe-events-calendar-month-mobile-events__mobile-event", func(e *colly.HTMLElement) {
		// limit到達チェック
		if stopped {
			return // 早期リターン
		}

		// モバイル版の日付を取得
		dateTimeText := e.ChildText(".tribe-events-calendar-month-mobile-events__mobile-event-datetime")

		// モバイル版のタイトルを取得
		titleText := e.ChildText(".tribe-events-calendar-month-mobile-events__mobile-event-title")

		// 見つけたイベント情報をコンソールに表示
		fmt.Printf("Found event: DateTime=%s, Title=%s\n", dateTimeText, titleText)

		// 日付の処理: 終日イベントの場合はdatetime属性から取得
		var formattedDate string
		if strings.TrimSpace(dateTimeText) == "終日" {
			// 終日イベントの場合、datetime属性から日付を取得
			datetimeAttr := e.ChildAttr("time", "datetime")
			formattedDate = datetimeAttr // 既にYYYY-MM-DD形式
			fmt.Printf("終日イベント検出: datetime=%s\n", datetimeAttr)
		} else {
			// 通常のイベントの場合、従来の処理
			datePart := strings.Split(dateTimeText, " /")[0]
			datePart = strings.TrimSpace(datePart)
			formattedDate = strings.ReplaceAll(datePart, ".", "-")
		}

		// 重複チェック: Repository層のメソッドを使用
		exists, err := s.eventRepo.ExistsByDateAndTitle(formattedDate, titleText)
		if err != nil {
			fmt.Printf("Error checking duplicate: %v\n", err)
			return
		}

		if exists {
			// 既に同じイベントがあるので、このイベントは保存しない
			fmt.Printf("スキップ: %s (%s) - すでに存在\n", titleText, formattedDate)
			return
		}

		// イベントの詳細説明を取得
		fullDescription := e.ChildText("div.tribe-events-single-event-description")

		newEvent := repository.Event{
			Date:        formattedDate,
			Title:       titleText,
			IsAttending: false,
			Details:     map[string]string{"description": fullDescription},
		}

		// Repository層のCreateメソッドを使用してDB保存
		createdEvent, err := s.eventRepo.Create(newEvent)
		if err != nil {
			fmt.Printf("Error creating event: %v\n", err)
			return
		}

		newEvents = append(newEvents, createdEvent)
		fmt.Printf("追加: %s (%s) - ID: %d\n", createdEvent.Title, createdEvent.Date, createdEvent.ID)

		// limit到達チェック
		if limit > 0 && len(newEvents) >= limit {
			stopped = true
			fmt.Printf("🛑 limit到達: %d件取得完了\n", limit)
		}
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
