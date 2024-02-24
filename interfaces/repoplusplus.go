package interfaces

import "gorm.io/gorm"

type Plusplus struct {
	gorm.Model `exhaustruct:"optional"`
	Name       string `gorm:"<-:create;uniqueIndex"`
	Value      int    `gorm:"<-;index"`
}

type PlusplusRepoInterface interface {
	Increment(name string, increment int) (int, error)
	FindTops(limit int) ([]Plusplus, error)
	FindFlops(limit int) ([]Plusplus, error)
}
