package entity

import (
	"database/sql"
)

type AbstractEntity struct {
	ID        int64 `gorm:"primaryKey"`
	CreatedBy sql.NullInt64
	UpdatedBy sql.NullInt64
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}
