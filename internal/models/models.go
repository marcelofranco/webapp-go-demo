package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
}

type Room struct {
	gorm.Model
	RoomName string
}

type Restriction struct {
	gorm.Model
	RestrictionName string
}

// Reservation holds reservation data
type Reservation struct {
	gorm.Model
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	Room      Room
	Processed int
}

type RoomRestriction struct {
	gorm.Model
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}

// MailData holds email message
type MailData struct {
	From     string
	To       string
	Subject  string
	Content  string
	Template string
}
