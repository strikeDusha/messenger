package database

import (
	"database/sql"
	"time"
)

type User struct {
	PublicName string `json:"pn"  binding:"required,min=2,max=30"`
	Username   string `json:"un"  binding:"required,min=5,max=20"` //later bind
	Email      string `json:"email"  binding:"required,email"`
	Password   string `json:"password"  binding:"required"` //later bind
}
type Message struct {
	Sender   string    `json:"sender"`
	Receiver string    `json:"receiver"`
	Text     string    `json:"text"`
	Sent     time.Time `json:"timestamp"`
}
type MsgMini struct {
	Text string    `json:"text"`
	Time time.Time `json:"timestamp"`
}
type Storage struct {
	db *sql.DB
}
