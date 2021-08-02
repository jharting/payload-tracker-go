package models

import (
	"time"
)

type PayloadStatuses struct {
	ID  uint `gorm:"primaryKey;not null;autoIncrement"`
	PayloadId uint `gorm:"not null"`
	ServiceId int32 `gorm:"not null"`
	SourceId int32 `gorm:"source_id"`
	StatusId int32 `gorm:"not null"`
	StatusMsg string `gorm:"type:varchar(100)"`
	Date time.Time `gorm:"primaryKey;not null`
	CreatedAt time.Time `gorm:"not null"`
}

type Payloads struct {
	Id uint `gorm:"primaryKey;not null"`
	RequestId string `gorm:"not null;type:varchar(100)"`
	Account string `gorm:"type:varchar(100)"`
	InventoryId string `gorm:"type:varchar(100)"`
	SystemId string `gorm:"type:varchar(100)"`
	CreatedAt time.Time `gorm:"not null"`
}

type Services struct {
	Id int32 `gorm:"primaryKey;not null"`
	Name string `gorm:"not null;type:varchar(100)"`
}

type Sources struct {
	Id int32 `gorm:"primaryKey;not null"`
	Name string `gorm:"not null;type:varchar(100)"`
}

type Statuses struct {
	Id int32 `gorm:"primaryKey;not null"`
	Name string `gorm:"not null;type:varchar(100)"`
}
