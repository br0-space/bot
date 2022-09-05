package interfaces

import "gorm.io/gorm"

type Plusplus struct {
	gorm.Model
	ChatID int64  `gorm:"<-:create;index:idx_plusplus_chat_id_user_id,unique"`
	Name   string `gorm:"<-:create;index:idx_plusplus_chat_id_user_id,unique"`
	Value  int    `gorm:"<-;index"`
}

type PlusplusRepoInterface interface {
	Increment(chatID int64, name string, increment int) (int, error)
	FindTops(chatID int64) ([]Plusplus, error)
	FindFlops(chatID int64) ([]Plusplus, error)
}
