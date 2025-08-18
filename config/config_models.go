package config

import (
	"time"

	"github.com/google/uuid"
)

type CreationTime struct {
	ID        *uuid.UUID `gorm:"type:text;default:uuid_generate_v4();primary_key"`
	CreatedAt *time.Time `gorm:"not null;type:timestamp;default:now()"`
	UpdatedAt *time.Time `gorm:"not null;type:timestamp;default:now()"`
}

type User struct {
	Name         string  `gorm:"type:varchar(100);not null"`
	UserName     string  `gorm:"not null;column:username;unique" json:"username"`
	PasswordHash string  `gorm:"not null;column:password" json:"password_hash,omitempty"`
	Role         *string `gorm:"type:varchar(50);default:'user';not null"`
	Verified     *bool   `gorm:"not null;default:false"`
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
	Title  string `gorm:"not null" json:"title"`
	Author string `gorm:"not null" json:"author"`
	CreationTime
}

type UserKeys struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	PGPKey string    `gorm:"type:string;column:pgpkey"`
	CreationTime
}

// The request for both register and login
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
