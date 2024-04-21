package main

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type BanRecord struct {
	ID      uint
	UserID  int
	BanType int
	BanDate time.Time
}

func main() {
	db, err := gorm.Open(sqlite.Open("ban_records.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var banRecords []BanRecord
	now := time.Now()
	err = db.Table("ban_records").
		Where("ban_date < ?", now).
		Order("ban_date DESC").
		Find(&banRecords).Error

	if err != nil {
		fmt.Println("Query failed:", err)
		return
	}

	fmt.Printf("%+v\n", banRecords)

	// closestBansの生成
	tempClosestBans := make(map[int]BanRecord)
	for _, record := range banRecords {
		// BanTypeが2のレコードはスキップ
		if record.BanType == 2 {
			continue
		}

		// 同じUserIDのレコードがまだマップにない、またはより古いban_dateを持つ場合にのみ、マップに追加/更新
		if existingRecord, exists := tempClosestBans[record.UserID]; !exists || existingRecord.BanDate.Before(record.BanDate) {
			tempClosestBans[record.UserID] = record
		}
	}

	closestBans := make(map[int]bool)
	for userID := range tempClosestBans {
		closestBans[userID] = true
	}
	fmt.Printf("%+v\n", closestBans)

	// 必要なUserIDのリスト
	requiredUserIDs := []uint32{3, 4, 5, 8, 10}

	// closestBansに含まれるUserIDをフィルタリングして除外
	remainingUserIDs := []uint32{}
	for _, id := range requiredUserIDs {
		if _, found := closestBans[int(id)]; !found {
			remainingUserIDs = append(remainingUserIDs, id)
		}
	}

	// 結果の表示
	fmt.Println("残ったUserID:", remainingUserIDs)
}
