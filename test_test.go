package main

import (
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	UserID       int `gorm:"primaryKey"`
	TargetUserID int `gorm:"primaryKey"`
	Name         string
	IsBanned     bool
	Timestamp    time.Time
}

type BanUser struct {
	UserID int `gorm:"primaryKey"`
	Reason string
}

const maxUsers = 5 // 仮の上限数

func manageUser2(tx *gorm.DB, checkUserID, checkTargetUserID int, isInsertOperation bool) error {
	var totalUsers int64
	var bannedUsers int64

	// 現在のユーザー数をカウント
	tx.Model(&User{}).Count(&totalUsers)

	// // BANユーザー数をカウント
	tx.Model(&BanUser{}).Count(&bannedUsers)

	bannedUsers = 0 // 仮の上限数のため、BANユーザー数は無視

	effectiveLimit := maxUsers - bannedUsers
	// effectiveLimit := int64(maxUsers)

	var matchedUser User
	result := tx.Where("user_id = ? AND target_user_id = ?", checkUserID, checkTargetUserID).First(&matchedUser)

	if isInsertOperation {
		if totalUsers >= effectiveLimit {
			return fmt.Errorf("cannot insert new user: limit reached (effective limit: %d)", effectiveLimit)
		}
		if result.Error == gorm.ErrRecordNotFound {
			newUser := User{UserID: checkUserID, TargetUserID: checkTargetUserID, Name: "New User", Timestamp: time.Now()}
			if err := tx.Create(&newUser).Error; err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}
			return nil // Successfully added new user
		}
		return fmt.Errorf("user already exists")
	} else {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("cannot delete: user does not exist")
		}
		if err := tx.Delete(&matchedUser).Error; err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}
		return nil // Successfully deleted user
	}
}

func TestManageUserInsertAndDeleteWithBans2(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	if err := db.AutoMigrate(&User{}, &BanUser{}); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// BANユーザーを追加
	db.Create(&BanUser{UserID: 2, Reason: "Violation"})
	db.Create(&BanUser{UserID: 4, Reason: "Inactive"})

	testCases := []struct {
		label         string
		userID        int
		targetUserID  int
		operation     bool // true for insert, false for delete
		expectError   bool
		expectedCount int64 // Expected number of users after the operation
	}{
		{"[テストパターン1] Insert user 1", 1, 40, true, false, 1},
		{"[テストパターン2] Insert BAN user 2", 2, 30, true, true, 1},
		{"[テストパターン3] Delete user 1", 1, 40, false, false, 0},
		{"[テストパターン4] Insert user 3", 3, 60, true, false, 1},
		{"[テストパターン5] Insert BAN user 4", 4, 20, true, true, 1},
		{"[テストパターン6] Insert user 5", 5, 50, true, false, 2},
		{"[テストパターン7] Insert user 6 - limit reached", 6, 70, true, true, 2},
		{"[テストパターン8] Delete user 5", 5, 50, false, false, 1},
		{"[テストパターン9] Delete non-existent user 6", 6, 70, false, true, 1},
	}

	for _, tc := range testCases {
		err := manageUser2(db, tc.userID, tc.targetUserID, tc.operation)
		if (err != nil && !tc.expectError) || (err == nil && tc.expectError) {
			t.Errorf("%s: Test failed for userID %d and targetUserID %d, operation %v: %v", tc.label, tc.userID, tc.targetUserID, tc.operation, err)
		} else {
			operationType := "inserted"
			if !tc.operation {
				operationType = "deleted"
			}
			if err == nil {
				t.Logf("%s: Successfully %s user with userID %d and targetUserID %d", tc.label, operationType, tc.userID, tc.targetUserID)
			} else {
				t.Logf("%s: Expected error on %s user with userID %d and targetUserID %d: %v", tc.label, operationType, tc.userID, tc.targetUserID, err)
			}
		}

		// Check the user count after each operation
		var userCount int64
		db.Model(&User{}).Count(&userCount)
		if userCount != tc.expectedCount {
			t.Errorf("%s: Expected %d users in the database, found %d, after operation on userID %d", tc.label, tc.expectedCount, userCount, tc.userID)
		}
	}
}
