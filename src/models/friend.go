package model

type FriendRequest struct {
	RequesterID uint `gorm:"primaryKey"`
	ReceiverID  uint `gorm:"primaryKey"`
	Requester   User `gorm:"foreignKey:RequesterID"`
	Receiver    User `gorm:"foreignKey:ReceiverID"`
}
