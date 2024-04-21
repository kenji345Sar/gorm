package main

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	UserID       int `gorm:"primaryKey;autoIncrement:false"`
	TargetUserID int `gorm:"primaryKey;autoIncrement:false"`
	Name         string
	Timestamp    time.Time
}

func main() {
	// Gorm DB接続の初期化（SQLiteを使用）
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// テーブルがなければ作成
	db.AutoMigrate(&User{})

	// サンプルデータの挿入
	users := []User{
		{UserID: 3, TargetUserID: 10, Name: "Alice", Timestamp: time.Date(2024, 4, 10, 23, 0, 0, 0, time.UTC)},
		{UserID: 1, TargetUserID: 20, Name: "Charlie", Timestamp: time.Date(2024, 4, 10, 22, 0, 0, 0, time.UTC)},
		{UserID: 2, TargetUserID: 30, Name: "Bob", Timestamp: time.Date(2024, 4, 10, 22, 30, 0, 0, time.UTC)},
		{UserID: 2, TargetUserID: 30, Name: "Dave", Timestamp: time.Date(2024, 4, 9, 21, 0, 0, 0, time.UTC)},
		{UserID: 1, TargetUserID: 40, Name: "Eve", Timestamp: time.Date(2024, 4, 11, 23, 0, 0, 0, time.UTC)},
		{UserID: 1, TargetUserID: 40, Name: "Eve5", Timestamp: time.Date(2024, 4, 15, 23, 0, 0, 0, time.UTC)},
	}
	db.Create(&users) // データベースにユーザーを挿入

	var userCount int64
	db.Model(&User{}).Count(&userCount)
	println(fmt.Sprintf("Inserted User: UserID: %d", userCount))

	// 挿入されたデータを確認
	var insertedUsers []User
	db.Find(&insertedUsers)
	for _, user := range insertedUsers {
		println(fmt.Sprintf("Inserted User: UserID: %d, TargetUserID: %d, Name: %s, Timestamp: %s",
			user.UserID, user.TargetUserID, user.Name, user.Timestamp))
	}
}
