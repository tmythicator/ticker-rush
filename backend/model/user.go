package model

import "time"

type User struct {
	ID           int64          `json:"id"`
	Email        string         `json:"email"`
	PasswordHash string         `json:"-"` // never send this to client
	Balance      float64        `json:"balance"`
	Portfolio    map[string]int `json:"portfolio"`
	CreatedAt    time.Time      `json:"created_at"`
}
