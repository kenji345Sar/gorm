package main

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// BanRecord テーブルを表す構造体
type BanRecord struct {
	ID      uint `gorm:"primaryKey"`
	UserID  int
	BanType int
	BanDate time.Time
}

func main() {
	// 指定したファイル名のSQLiteデータベースに接続
	db, err := gorm.Open(sqlite.Open("ban_records.db"), &gorm.Config{})
	if err != nil {
		panic("データベースへの接続に失敗しました")
	}

	// `ban_records` テーブルを自動で作成
	err = db.AutoMigrate(&BanRecord{})
	if err != nil {
		panic("テーブルのマイグレーションに失敗しました")
	}

	// テーブル内の既存データを削除
	err = db.Where("1 = 1").Delete(&BanRecord{}).Error
	if err != nil {
		panic("データの削除に失敗しました")
	}

	// 新しいデータを挿入
	banRecords := []BanRecord{
		{ID: 1, UserID: 30, BanType: 1, BanDate: time.Date(2024, 4, 5, 15, 11, 13, 0, time.UTC)},
		{ID: 2, UserID: 20, BanType: 1, BanDate: time.Date(2024, 4, 3, 15, 11, 13, 0, time.UTC)},
		{ID: 3, UserID: 5, BanType: 1, BanDate: time.Date(2024, 3, 28, 15, 11, 13, 0, time.UTC)},
		{ID: 4, UserID: 5, BanType: 1, BanDate: time.Date(2024, 3, 29, 15, 11, 13, 0, time.UTC)},
		{ID: 5, UserID: 5, BanType: 1, BanDate: time.Date(2024, 4, 1, 15, 11, 13, 0, time.UTC)},
		{ID: 6, UserID: 8, BanType: 1, BanDate: time.Date(2024, 3, 28, 15, 11, 13, 0, time.UTC)},
		{ID: 7, UserID: 8, BanType: 1, BanDate: time.Date(2024, 3, 29, 15, 11, 13, 0, time.UTC)},
		{ID: 8, UserID: 8, BanType: 2, BanDate: time.Date(2024, 4, 1, 15, 11, 13, 0, time.UTC)},
	}

	for _, record := range banRecords {
		result := db.Create(&record) // レコードを挿入
		if result.Error != nil {
			fmt.Println("データ挿入時にエラーが発生しました:", result.Error)
			return
		}
	}

	fmt.Println("データの挿入が完了しました")
}
