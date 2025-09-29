package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"kadai7_database-integration/repository"
	"kadai7_database-integration/service"
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

	// Repository層の初期化
	repo, err := repository.NewEventRepository()
	if err != nil {
		return fmt.Errorf("Repository初期化エラー: %v", err)
	}
	defer repo.Close()

	// Service層の初期化
	eventService := service.NewEventService(repo)

	// スクレイピング実行
	events, err := eventService.ScrapeAndSaveEvents(limit)
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
