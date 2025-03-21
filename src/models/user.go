package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string  `json:"username" gorm:"unique;not null"`
	Email    string  `json:"email" gorm:"unique;not null"`
	Phone    *string `json:"phone"`
	Password []byte  `json:"-"`
	Session  *string `json:"-"`
	Friends  []*User `json:"friends" gorm:"many2many:user_friends;joinForeignKey:UserID;JoinReferences:FriendID"`
}
