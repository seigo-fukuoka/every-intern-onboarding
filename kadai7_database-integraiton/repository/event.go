package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	migrate "github.com/rubenv/sql-migrate"
)

// Event - イベント構造体（shared/models.goから移行）
type Event struct {
	ID          int               `json:"id"`          // 自動生成ID
	Date        string            `json:"date"`        // 日付 (YYYY-MM-DD)
	Title       string            `json:"title"`       // イベント名
	IsAttending bool              `json:"isAttending"` // 参加予定フラグ
	Details     map[string]string `json:"details"`     // 追加情報（JSON形式）
}

// EventRepository - イベントのデータアクセス層
type EventRepository struct {
	db *sql.DB
}

// NewEventRepository - EventRepositoryのコンストラクタ（DB接続 + 初期化）
func NewEventRepository() (*EventRepository, error) {
	// DB接続（shared/database.goのConnectDBから移行）
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/events_db?parseTime=true")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("データベース接続成功！")

	repo := &EventRepository{db: db}

	// マイグレーション実行（shared/database.goのInitDBから移行）
	if err := repo.RunMigrations(); err != nil {
		return nil, err
	}

	return repo, nil
}

// RunMigrations - データベースの構造変更を実行（shared/database.goから移行）
func (r *EventRepository) RunMigrations() error {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations/",
	}
	n, err := migrate.Exec(r.db, "mysql", migrations, migrate.Up)
	if err != nil {
		return err
	}
	fmt.Printf("Applied %d migrations!\n", n)
	return nil
}

// GetAll - 全イベント取得（main.goのDB操作から移行）
func (r *EventRepository) GetAll() ([]Event, error) {
	query := `SELECT id, date, title, is_attending, details FROM events ORDER BY date ASC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]Event, 0)
	for rows.Next() {
		var event Event
		var detailsJSON string

		err := rows.Scan(&event.ID, &event.Date, &event.Title, &event.IsAttending, &detailsJSON)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(detailsJSON), &event.Details); err != nil {
			event.Details = make(map[string]string)
		}

		events = append(events, event)
	}

	return events, nil
}

// ExistsByDateAndTitle - 重複チェック（shared/scraper.goから移行）
func (r *EventRepository) ExistsByDateAndTitle(date, title string) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM events WHERE date = ? AND title = ?", date, title).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Create - イベント作成（shared/scraper.goから移行）
func (r *EventRepository) Create(event Event) (Event, error) {
	detailsJSON, err := json.Marshal(event.Details)
	if err != nil {
		return event, err
	}

	result, err := r.db.Exec("INSERT INTO events (date, title, is_attending, details) VALUES (?, ?, ?, ?)",
		event.Date, event.Title, event.IsAttending, string(detailsJSON))
	if err != nil {
		return event, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return event, err
	}
	event.ID = int(lastInsertID)

	return event, nil
}

// Close - DB接続を閉じる
func (r *EventRepository) Close() error {
	return r.db.Close()
}
