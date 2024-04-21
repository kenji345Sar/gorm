package main

import (
	"fmt"
	"sort"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User モデルの定義: user_idカラムを追加
type User struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"index"` // user_idはインデックスを付与
	Name      string
	Timestamp time.Time
}

func main() {
	// SQLiteデータベースに接続（他のデータベースでの接続も同様に行えます）
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("データベースへの接続に失敗しました")
	}

	// データベースの自動マイグレーション（テーブルがない場合は作成します）
	db.AutoMigrate(&User{})

	// テーブルの内容をクリア
	db.Exec("DELETE FROM users") // SQLiteの場合、これで全データを削除

	// テストデータの挿入
	users := []User{
		{UserID: 3, Name: "Alice", Timestamp: time.Date(2024, 4, 10, 23, 0, 0, 0, time.UTC)},
		{UserID: 1, Name: "Charlie", Timestamp: time.Date(2024, 4, 10, 22, 0, 0, 0, time.UTC)},
		{UserID: 2, Name: "Bob", Timestamp: time.Date(2024, 4, 10, 22, 30, 0, 0, time.UTC)},
		{UserID: 2, Name: "Dave", Timestamp: time.Date(2024, 4, 9, 21, 0, 0, 0, time.UTC)},
		{UserID: 1, Name: "Eve", Timestamp: time.Date(2024, 4, 11, 23, 0, 0, 0, time.UTC)},
		{UserID: 1, Name: "Eve5", Timestamp: time.Date(2024, 4, 15, 23, 0, 0, 0, time.UTC)},
	}
	for _, user := range users {
		db.Create(&user) // 各ユーザーをデータベースに挿入
	}

	// ユーザーデータの取得
	var usersRetrieved []User
	result := db.Find(&usersRetrieved)
	if result.Error != nil {
		fmt.Println("データの取得中にエラーが発生しました:", result.Error)
		return
	}

	// 取得したデータをまず user_id で昇順、次に Timestamp で降順に並び替え
	sort.Slice(usersRetrieved, func(i, j int) bool {
		if usersRetrieved[i].UserID != usersRetrieved[j].UserID {
			return usersRetrieved[i].UserID < usersRetrieved[j].UserID
		}
		return usersRetrieved[i].Timestamp.After(usersRetrieved[j].Timestamp)
	})

	// 並び替えたデータの出力
	for _, user := range usersRetrieved {
		fmt.Printf("UserID: %d, Name: %s, Timestamp: %s\n", user.UserID, user.Name, user.Timestamp)
	}
}
