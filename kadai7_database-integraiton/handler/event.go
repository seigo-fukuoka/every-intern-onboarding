package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"kadai7_database-integration/service"
)

// EventHandler - イベントのプレゼンテーション層（HTTP処理）
type EventHandler struct {
	eventService service.EventServiceInterface
}

// NewEventHandler - EventHandlerのコンストラクタ
func NewEventHandler(eventService service.EventServiceInterface) *EventHandler {
	return &EventHandler{
		eventService: eventService,
	}
}

// GetEvents - 全イベント取得API（main.goのGET /eventsから移行）
func (h *EventHandler) GetEvents(c echo.Context) error {
	// 「h.eventService」のGetAllEventsメソッドを呼び出し、Service層に処理を委譲
	events, err := h.eventService.GetAllEvents()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	// 成功レスポンス
	return c.JSON(http.StatusOK, events)
}

// ScrapeEvents - スクレイピング実行API（main.goのGET /scrape/all-eventsから移行）
func (h *EventHandler) ScrapeEvents(c echo.Context) error {
	// 「h.eventService」のScrapeAndSaveEventsメソッドを呼び出し、Service層のスクレイピング処理を呼び出し（limit=0で無制限）
	events, err := h.eventService.ScrapeAndSaveEvents(0)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	// 成功レスポンス
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "スクレイピング完了！",
		"count":   len(events),
		"events":  events,
	})
}
