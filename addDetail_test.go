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
	Timestamp    time.Time
}

// manageUser checks for a user and deletes if exists, or creates if not.
func manageUser(tx *gorm.DB, checkUserID, checkTargetUserID int) error {
	var matchedUser User
	result := tx.Where("user_id = ? AND target_user_id = ?", checkUserID, checkTargetUserID).First(&matchedUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Record not found, create a new user
			newUser := User{UserID: checkUserID, TargetUserID: checkTargetUserID, Name: "New User", Timestamp: time.Now()}
			if err := tx.Create(&newUser).Error; err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}
			fmt.Println("No existing user found. Added new user:", newUser)
			return nil
		}
		return fmt.Errorf("database error: %w", result.Error)
	}
	// User exists, delete it
	if err := tx.Delete(&matchedUser).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	fmt.Println("Existing user found. Deleted user:", matchedUser)
	return nil
}

// TestManageUser tests the manageUser function
func TestManageUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	checkUserID := 1
	checkTargetUserID := 40

	// Test 1: Ensure it inserts correctly when no user exists
	if err := manageUser(db, checkUserID, checkTargetUserID); err != nil {
		t.Fatalf("manageUser failed on insert: %v", err)
	}

	// Verify insertion
	var count int64
	db.Model(&User{}).Where("user_id = ? AND target_user_id = ?", checkUserID, checkTargetUserID).Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 user, got %d users", count)
	}

	// Test 2: Ensure it deletes correctly when the user exists
	if err := manageUser(db, checkUserID, checkTargetUserID); err != nil {
		t.Fatalf("manageUser failed on delete: %v", err)
	}

	// Verify deletion
	db.Model(&User{}).Where("user_id = ? AND target_user_id = ?", checkUserID, checkTargetUserID).Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 users, got %d users", count)
	}
}
