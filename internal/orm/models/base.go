package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// BaseModel defines the common columns that all db structs should hold, usually
// db structs based on this have no soft delete
type BaseModel struct {
	// Default values for PostgreSQL, change it for other DBMS
	ID        uuid.UUID  `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt *time.Time `gorm:"index;not null;default:current_timestamp"`
	UpdatedAt *time.Time `gorm:"index;not null;default:current_timestamp"`
}

// BaseModelSoftDelete defines the common columns that all db structs should
// hold, usually. This struct also defines the fields for GORM triggers to
// detect the entity should soft delete
type BaseModelSoftDelete struct {
	BaseModel
	DeletedAt *time.Time `gorm:"index"`
}

// BaseModelSeq defines the common columns that all db structs should hold, with
// an INT key
type BaseModelSeq struct {
	// Default values for PostgreSQL, change it for other DBMS
	ID        int        `gorm:"primary_key,auto_increment"`
	CreatedAt *time.Time `gorm:"index;not null;default:current_timestamp"`
	UpdatedAt *time.Time `gorm:"index;not null;default:current_timestamp"`
}

// BaseModelSeqSoftDelete defines the common columns that all db structs should
// hold, usually. This struct also defines the fields for GORM triggers to
// detect the entity should soft delete
type BaseModelSeqSoftDelete struct {
	BaseModelSeq
	DeletedAt *time.Time `gorm:"index"`
}
