package model

import (
	"time"
)

type EventLog struct {
	ID        uint      `gorm:"primaryKey"`
	Service   string    `gorm:"size:100;not null"` // e.g., "consumer", "mqtt_subscriber"
	Message   string    `gorm:"type:text;not null"`
	Context   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
