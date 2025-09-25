package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"kadai7_database-integraiton/shared"
)

func main() {
	// CLI アプリケーション設定
	// スクレイピング機能のみを実装
	app := &cli.App{
		Name:  "schedule-scraper",
		Usage: "スケジュールスクレイピング",
		Commands: []*cli.Command{
			{
				Name:  "scrape",
				Usage: "イベントをスクレイピング",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "limit",
						Value: 0,
						Usage: "取得する件数の上限(0は無制限)",
					},
				},
				Action: func(c *cli.Context) error {
					limit := c.Int("limit")
					return runScrapingWithLimit(limit)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)

	}
}

func runScrapingWithLimit(limit int) error {
	fmt.Println("🚀 イベントスクレイピングを開始...")

	// データベース接続
	db, err := shared.ConnectDB()
	if err != nil {
		return fmt.Errorf("データベース接続エラー: %v", err)
	}
	defer db.Close()

	// マイグレーション実行
	err = shared.InitDB(db)
	if err != nil {
		return fmt.Errorf("マイグレーションエラー: %v", err)
	}

	// スクレイピング実行
	events, err := shared.ScrapeAllEvent(db, limit)
	if err != nil {
		return fmt.Errorf("スクレイピングエラー: %v", err)
	}

	fmt.Printf("✅ スクレイピング完了！ %d件のイベントを取得しました\n", len(events))

	// 取得したイベントを表示
	for i, event := range events {
		fmt.Printf("%d. %s - %s\n", i+1, event.Date, event.Title)
	}

	return nil
}
