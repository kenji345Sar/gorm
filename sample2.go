package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User モデル
type User struct {
	gorm.Model
	Name string
}

// シャードされたデータベースへの接続を管理する
var shardMap = map[int]*gorm.DB{}

func init() {
	// シャードデータベースの初期化
	for i := 1; i <= 2; i++ {
		db, err := gorm.Open(sqlite.Open(fmt.Sprintf("test_shard_%d.db", i)), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		db.AutoMigrate(&User{})
		shardMap[i] = db
	}
}

// ユーザーIDに基づいてシャードを選択
func getShard(userID uint) *gorm.DB {
	// ここでは単純化のため、IDが奇数の場合はシャード1、偶数の場合はシャード2を選択
	if userID%2 == 0 {
		return shardMap[2]
	}
	return shardMap[1]
}

func main() {
	// ユーザーの作成
	user1 := User{Name: "Alice"}
	user2 := User{Name: "Bob"}

	// ユーザーを適切なシャードに挿入
	shard1 := getShard(1)
	shard2 := getShard(2)

	shard1.Create(&user1)
	fmt.Println("Inserted User ID:", user1.ID)
	var user User
	shard1.First(&user, user1.ID) // 修正: ユーザー1の情報をシャード1から取得
	fmt.Println(user)

	shard2.Create(&user2)
	fmt.Println("Inserted User ID:", user2.ID)
	shard2.First(&user, user2.ID) // 修正: ユーザー2の情報をシャード2から取得
	fmt.Println(user)
}
