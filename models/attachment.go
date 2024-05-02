package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Attachment struct {
	gorm.Model
	ID        string `gorm:"primaryKey" json:"id"`
	OwnerID   string `gorm:"not null" json:"owner_id"`
	Path      string `gorm:"not null" json:"path"`
	Extension string `gorm:"not null" json:"Extension"`
}

func (a *Attachment) BeforeCreate(*gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return
}
