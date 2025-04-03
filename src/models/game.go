package model

import "time"

type Game struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `json:"-"`
	User          User      `gorm:"foreignKey:UserID;references:ID" json:"user"`
	OpponentID    uint      `json:"-"`
	Opponent      User      `gorm:"foreignKey:OpponentID;references:ID" json:"opponent"`
	Board         string    `json:"board"`
	StatusID      uint      `json:"-"`
	Status        Status    `json:"status"`
	WinnerID      *uint     `json:"-"`
	Winner        *User     `gorm:"foreignKey:WinnerID;references:ID" json:"winner"`
	WhitePlayerID *uint     `json:"white_player_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Status struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Status    string    `json:"status" gorm:"unique;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	StatusPending = "PENDING"
	StatusPlaying = "PLAYING"
	StatusEnd     = "END"
)
