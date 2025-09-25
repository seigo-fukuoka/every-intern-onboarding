package shared

type Event struct {
	ID          int               `json:"id"`          // 自動生成ID
	Date        string            `json:"date"`        // 日付 (YYYY-MM-DD)
	Title       string            `json:"title"`       // イベント名
	IsAttending bool              `json:"isAttending"` // 参加予定フラグ
	Details     map[string]string `json:"details"`     // 追加情報（JSON形式）
}
