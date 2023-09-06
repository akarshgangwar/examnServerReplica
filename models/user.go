package models

import (
	"time"
)

// I will write the model here for user which will have the following fields:
// Email - string
// Password - string
// AccessToken - string
// RefreshToken - string
// TokenExpirationTime - time.Time
type User struct {
	Email               string     `gorm:"primaryKey"`
	Password            string     `json:"password"`
	AccessToken         string     `json:"access_token"`
	RefreshToken        string     `json:"refresh_token"`
	TokenExpirationTime time.Time `json:"token_expiration_time"`
}

func (s *User) TableName() string {
	return "user"
}
