package main

import (
	"errors"
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
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		fmt.Printf("Failed to migrate database: %v\n", err)
		return
	}

	if err := executeTransaction(db); err != nil {
		fmt.Printf("Transaction failed: %v\n", err)
	} else {
		fmt.Println("Transaction succeeded.")
	}

	showUsers(db)
}

func executeTransaction(db *gorm.DB) error {
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// トランザクションが失敗した場合にのみRollbackを実行
	defer func() {
		if tx.Error != nil {
			tx.Rollback()
		}
	}()

	checkUserID := 1
	checkTargetUserID := 40

	var matchedUser User

	result := tx.Where("user_id = ? AND target_user_id = ?", checkUserID, checkTargetUserID).First(&matchedUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// レコードが見つからなかった場合の処理
			newUser := User{UserID: checkUserID, TargetUserID: checkTargetUserID, Name: "New User", Timestamp: time.Now()}
			if err := tx.Create(&newUser).Error; err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}
			fmt.Println("No existing user found. Added new user:", newUser)
		} else {
			// その他のエラーの場合の処理
			return fmt.Errorf("database error: %w", result.Error)
		}
	} else {
		// エラーが発生しなかった場合（レコードが正常に見つかった場合）
		if err := tx.Delete(&matchedUser).Error; err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}
		fmt.Println("Existing user found. Deleted user:", matchedUser)
	}

	// エラーがなければCommitを実行
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func showUsers(db *gorm.DB) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		fmt.Printf("Failed to list users: %v\n", err)
		return
	}
	for _, user := range users {
		fmt.Printf("UserID: %d, TargetUserID: %d, Name: %s, Timestamp: %s\n",
			user.UserID, user.TargetUserID, user.Name, user.Timestamp)
	}
}
