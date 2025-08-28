package xorm

import (
	"time"
)

type BaseTimeModel struct {
	DeletedAt time.Time `xorm:"deleted index DATETIME(6)"`
	CreatedAt time.Time `xorm:"created not null DATETIME(6)"`
	CreatedBy string    `xorm:"not null VARCHAR(32)"`
	UpdatedAt time.Time `xorm:"updated not null DATETIME(6)"`
	UpdatedBy string    `xorm:"not null VARCHAR(32)"`
}

type BaseTimeModelWithoutSoftDelete struct {
	CreatedAt time.Time `xorm:"created not null DATETIME(6)"`
	CreatedBy string    `xorm:"not null VARCHAR(32)"`
	UpdatedAt time.Time `xorm:"updated not null DATETIME(6)"`
	UpdatedBy string    `xorm:"not null VARCHAR(32)"`
}
