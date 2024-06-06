package dbp

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string `gorm:"unique"`
	Email          string `gorm:"unique"`
	Password       string
	Image          string
	Biography      string
	Gender         string
	DesiredGender  string
	Sport          string
	SecondarySport string
	DateOfBirth    time.Time
	City           string
}

type Strike struct {
	gorm.Model
	UserAID int
	UserA   User
	UserBID int
	UserB   User
}

type Miss struct {
	gorm.Model
	UserAID int
	UserA   User
	UserBID int
	UserB   User
}

type Match struct {
	gorm.Model
	UserAID    int
	UserA      User
	UserBID    int
	UserB      User
	RoomName   string
	Historical string
}

type Stat struct {
	ID         uint `gorm:"primarykey"`
	Name       string
	Catégories string
	Image      string
}
