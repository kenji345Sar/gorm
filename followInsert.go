package main

import (
	"fmt"

	"gorm.io/driver/sqlite" // SQLiteドライバをインポート
	"gorm.io/gorm"
)

// Follower テーブルを表す構造体
type Follower struct {
	Id         int `gorm:"primaryKey"`
	UserId     int
	FollowerId int
	BanFlg     int
}

func main() {
	// 指定したファイル名のSQLiteデータベースに接続
	db, err := gorm.Open(sqlite.Open("test_shard1.db"), &gorm.Config{})
	if err != nil {
		panic("データベースへの接続に失敗しました")
	}

	// `followers` テーブルを自動で作成
	err = db.AutoMigrate(&Follower{})
	if err != nil {
		panic("テーブルのマイグレーションに失敗しました")
	}

	// テーブル内の既存データを削除
	err = db.Where("1 = 1").Delete(&Follower{}).Error
	if err != nil {
		panic("データの削除に失敗しました")
	}

	// 新しいデータを挿入
	followers := []Follower{
		{UserId: 21, FollowerId: 5, BanFlg: 0},
		{UserId: 5, FollowerId: 21, BanFlg: 1},
		{UserId: 21, FollowerId: 30, BanFlg: 0},
		{UserId: 21, FollowerId: 35, BanFlg: 0},
		{UserId: 10, FollowerId: 21, BanFlg: 0},
		{UserId: 15, FollowerId: 21, BanFlg: 1},
		{UserId: 10, FollowerId: 15, BanFlg: 0}, // 新しいデータ
		{UserId: 11, FollowerId: 14, BanFlg: 0}, // 新しいデータ
		{UserId: 21, FollowerId: 11, BanFlg: 0}, // 新しいデータ
		{UserId: 11, FollowerId: 21, BanFlg: 0}, // 新しいデータ
		{UserId: 21, FollowerId: 18, BanFlg: 1}, // 新しいデータ
		{UserId: 18, FollowerId: 21, BanFlg: 0}, // 新しいデータ
	}

	for _, follower := range followers {
		result := db.Create(&follower) // レコードを挿入
		if result.Error != nil {
			fmt.Println("データ挿入時にエラーが発生しました:", result.Error)
			return
		}
	}

	fmt.Println("データの挿入が完了しました")
}
