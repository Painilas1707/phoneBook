package StructUser

import "time"

type UserPhoneBook struct {
	ID          int       `json:"id" gorm:"primary_key"`
	ContactFIO  string    `json:"contact_fio" validate:"required"`
	BirthDate   time.Time `json:"data_birth"`
	PhoneNumber string    `json:"phone_number" validate:"required,numeric" gorm:"index"`
	Email       string    `json:"email" validate:"email"`
	TimeCreate  time.Time `json:"TimeCreate"`
}
