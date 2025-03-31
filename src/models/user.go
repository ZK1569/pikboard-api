package model

import "time"

type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Image     string    `json:"image" gorm:"image"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Phone     *string   `json:"phone"`
	Password  []byte    `json:"-"`
	Session   *string   `json:"-"`
	Friends   []*User   `json:"friends" gorm:"many2many:user_friends;joinForeignKey:UserID;JoinReferences:FriendID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
