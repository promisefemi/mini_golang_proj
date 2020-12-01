package model

import (
	"errors"

	"github.com/twinj/uuid"
	"gorm.io/gorm"
)

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	p.UUID = uuid.New([]byte("Abunchnumbers")).String()
	if p.UUID == "" {
		return errors.New("can't save invalid data")
	}
	return nil
}
