package model

import (
	"time"

	"github.com/google/uuid"
)

type GenerationStatus string

const (
	GenerationStatusPending  GenerationStatus = "pending"
	GenerationStatusComplete GenerationStatus = "complete"
	GenerationStatusFailed   GenerationStatus = "failed"
)

type Generation struct {
	ID        uuid.UUID        `gorm:"type:uuid;primaryKey"                     json:"id"`
	PlateID   uuid.UUID        `gorm:"type:uuid;not null;index"                 json:"plate_id"`
	AccountID *uuid.UUID       `gorm:"type:uuid;index"                          json:"account_id,omitempty"`
	Status    GenerationStatus `gorm:"type:varchar(20);default:pending;not null" json:"status"`
	Values    []byte           `gorm:"type:jsonb"                               json:"values,omitempty"`
	Error     *string          `gorm:"type:text"                                json:"error,omitempty"`
	CreatedAt time.Time        `gorm:"autoCreateTime"                           json:"created_at"`
	UpdatedAt time.Time        `gorm:"autoUpdateTime"                           json:"updated_at"`

	Plate   *Plate   `gorm:"foreignKey:PlateID"   json:"plate,omitempty"`
	Account *Account `gorm:"foreignKey:AccountID" json:"account,omitempty"`
}

func (Generation) TableName() string {
	return "generation"
}
