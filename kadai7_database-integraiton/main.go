// きゅるりんってしてみてスケジュールアプリ - データベース統合版
// 機能: Webページから情報取得 + データベース保存 + API提供 + ブラウザ画面表示
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"kadai7_database-integraiton/shared"
)

// main - アプリケーションのエントリーポイント
func main() {
	if err := startServer(); err != nil {
		fmt.Printf("サーバー起動エラー: %v\n", err)
	}
}

func startServer() error {

	// Webサーバー初期化（Echo = Go用のWebフレームワーク）
	e := echo.New()

	// データベース接続
	db, err := shared.ConnectDB()
	if err != nil {
		return err
	}
	// マイグレーション実行（テーブル作成・更新）
	err = shared.InitDB(db)
	if err != nil {
		return err
	}

	// CORS設定（ブラウザからのアクセス許可設定）
	// React（ポート5173-5175）からGo（ポート1323）へのアクセスを許可
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:5175"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// API エンドポイント設定

	// GET /events - 全イベント取得
	e.GET("/events", func(c echo.Context) error {
		query := `SELECT id, date, title, is_attending, details FROM events ORDER BY date ASC`
		rows, err := db.Query(query)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}
		defer rows.Close()

		events := make([]shared.Event, 0) // 空スライスを明示的に初期化
		for rows.Next() {
			var event shared.Event
			var detailsJSON string

			err := rows.Scan(&event.ID, &event.Date, &event.Title, &event.IsAttending, &detailsJSON)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": err.Error(),
				})
			}

			if err := json.Unmarshal([]byte(detailsJSON), &event.Details); err != nil {
				event.Details = make(map[string]string)
			}

			events = append(events, event)
		}

		return c.JSON(http.StatusOK, events)
	})

	// GET /scrape/all-events - 全イベントスクレイピング実行
	e.GET("/scrape/all-events", func(c echo.Context) error {
		events, err := shared.ScrapeAllEvent(db, 0) // API側はlimit無制限
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "スクレイピング完了！",
			"count":   len(events),
			"events":  events,
		})
	})

	// サーバー起動（ポート1323で待機）
	fmt.Println("サーバーを起動中... http://localhost:1323")
	return e.Start(":1323")
}
