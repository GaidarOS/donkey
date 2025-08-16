package config

import (
	"time"

	"github.com/google/uuid"
)

type CreationTime struct {
	CreatedAt time.Time `gorm:"type:timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp"`
}

type User struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserName     string `gorm:"not null;column:username;unique" json:"username"`
	PasswordHash string `gorm:"not null;column:password" json:"password_hash,omitempty"`
	CreationTime
}

type Config struct {
	Dir                string          `json:"dir" default:"uploads"`
	Port               string          `json:"port" default:"8080"`
	Depth              string          `json:"depth" default:"3"`
	ConfFile           string          `json:"confFile" default:"./config.json"`
	AllowedHeaderTypes map[string]bool `json:"allowedHeaderTypes"`
	Users              []User          `json:"users"`
}

// book.go
type Book struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Title  string `gorm:"not null" json:"title"`
	Author string `gorm:"not null" json:"author"`
	CreationTime
}

type UserKeys struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	PGPKey string    `gorm:"type:string;column:pgpkey"`
	CreationTime
}

// The request Dto for both register and login
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
