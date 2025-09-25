package shared

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	migrate "github.com/rubenv/sql-migrate"
)

// ConnectDB - MySQLデータベースへの接続
func ConnectDB() (*sql.DB, error) {
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

// InitDB - データベース初期化（マイグレーション実行）
func InitDB(db *sql.DB) error {
	return RunMigrations(db)
}

// runMigrations - データベースの構造変更を実行
// migrations/フォルダ内の.sqlファイル（テーブル作成命令書）を順番に適用
func RunMigrations(db *sql.DB) error {
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
