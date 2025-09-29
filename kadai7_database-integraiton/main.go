// きゅるりんってしてみてスケジュールアプリ - レイヤードアーキテクチャ版
// 機能: ルーター専用（依存性注入 + ルーティング設定）
package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"kadai7_database-integration/handler"
	"kadai7_database-integration/repository"
	"kadai7_database-integration/service"
)

// main - アプリケーションのエントリーポイント
func main() {
	if err := startServer(); err != nil {
		fmt.Printf("サーバー起動エラー: %v\n", err)
	}
}

func startServer() error {
	// 依存性注入（Dependency Injection）
	// Repository層 → Service層 → Handler層の順序で初期化

	// 1. Repository層の初期化（DB接続 + マイグレーション）
	eventRepo, err := repository.NewEventRepository()
	if err != nil {
		return fmt.Errorf("Repository初期化エラー: %v", err)
	}
	defer eventRepo.Close() // サーバー終了時にDB接続を閉じる

	// 2. Service層の初期化（Repository層を注入）
	eventService := service.NewEventService(eventRepo)

	// 3. Handler層の初期化（Service層を注入）
	eventHandler := handler.NewEventHandler(eventService)

	// Webサーバー初期化（Echo = Go用のWebフレームワーク）
	e := echo.New()

	// CORS設定（ブラウザからのアクセス許可設定）
	// React（ポート5173-5175）からGo（ポート1323）へのアクセスを許可
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:5175"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// ルーティング設定（Handler層のメソッドを指定）
	e.GET("/events", eventHandler.GetEvents)               // 全イベント取得
	e.GET("/scrape/all-events", eventHandler.ScrapeEvents) // スクレイピング実行

	// サーバー起動（ポート1323で待機）
	fmt.Println("サーバーを起動中... http://localhost:1323")
	return e.Start(":1323")
}
