package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Password represents the passwords table in the database.
type Password struct {
	ID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null"`
	Hash   string    `gorm:"type:varchar(255);not null"`
	gorm.Model
}

// BeforeSave is a GORM hook that can be used to set the updated_at field before saving.
func (p *Password) BeforeSave(tx *gorm.DB) (err error) {
	p.UpdatedAt = time.Now()
	return
}
