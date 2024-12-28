package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email         string         `gorm:"type:varchar(255);unique;not null"`
	FullName      string         `gorm:"type:varchar(255)"`
	Role          string         `gorm:"type:varchar(255);not null;default:user"`
	EmailVerified bool           `gorm:"not null;default:false"`
	Bio           string         `gorm:"type:text"`
	AvatarURL     string         `gorm:"type:text"`
	Metadata      datatypes.JSON `gorm:"type:jsonb"`
	gorm.Model
}

// BeforeCreate will set a UUID rather than numeric ID.
func (user *User) BeforeCreate(tx *gorm.DB) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return nil
}
