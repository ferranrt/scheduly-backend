package dbmodels

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Center struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string
	OwnerID   uuid.UUID `gorm:"type:uuid;not null"`
	Owner     User      `gorm:"foreignKey:OwnerID;references:ID"`
}

func (c *Center) TableName() string {
	return "centers"
}

func (c *Center) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return
}

func (c *Center) AfterUpdate(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now()
	return
}
