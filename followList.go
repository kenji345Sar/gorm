package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Follower struct {
	Id         int `gorm:"primaryKey"`
	UserId     int
	FollowerId int
	BanFlg     int
}

func main() {
	db, err := gorm.Open(sqlite.Open("test_shard1.db"), &gorm.Config{})
	if err != nil {
		panic("データベースへの接続に失敗しました")
	}

	var followers []Follower
	// ログインユーザーに関連するフォロー関係を取得
	if err := db.Where("user_id = ? OR follower_id = ?", 21, 21).Find(&followers).Error; err != nil {
		panic("フォロー関係の取得に失敗しました")
	}

	// フォローリストとフォロワーリストを生成して出力
	createLists(followers)
}

func createLists(followers []Follower) {
	const loginUser = 21

	following := make(map[int]bool)
	followedBy := make(map[int]bool)
	bannedUsers := make(map[int]bool)

	for _, f := range followers {
		// ログインユーザーがフォローしているユーザー
		if f.UserId == loginUser {
			following[f.FollowerId] = true
			// ban_flg=1であれば記録
			if f.BanFlg == 1 {
				bannedUsers[f.FollowerId] = true
			}
		}
		// ログインユーザーをフォローしているユーザー
		if f.FollowerId == loginUser {
			followedBy[f.UserId] = true
			// ban_flg=1であれば記録
			if f.BanFlg == 1 {
				bannedUsers[f.UserId] = true
			}
		}
	}

	// フォローリストを生成し、ban_flg=1であるユーザーを除外
	finalFollowing := []int{}
	for userId := range following {
		if !bannedUsers[userId] { // ban_flg=1であるユーザーを除外
			finalFollowing = append(finalFollowing, userId)
		}
	}

	// フォロワーリストを生成し、ban_flg=1であるユーザーを除外
	finalFollowedBy := []int{}
	for userId := range followedBy {
		if !bannedUsers[userId] { // ban_flg=1であるユーザーを除外
			finalFollowedBy = append(finalFollowedBy, userId)
		}
	}

	fmt.Println("フォローリスト:", finalFollowing)
	fmt.Println("フォロワーリスト:", finalFollowedBy)
}
