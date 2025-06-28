package schema

import (
	"fp-kpl/domain/identity"
	"fp-kpl/domain/shared"
	"fp-kpl/domain/table"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Table struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();column:id"`
	TableNumber string         `gorm:"type:varchar(255);unique;not null;column:table_number"`
	CreatedAt   time.Time      `gorm:"type:timestamp with time zone;column:created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamp with time zone;column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"type:timestamp with time zone;column:deleted_at"`

	Transactions []Transaction `gorm:"foreignKey:TableID"`
}

func TableEntityToSchema(entity table.Table) Table {
	var deletedAtTime time.Time
	if entity.DeletedAt != nil {
		deletedAtTime = *entity.DeletedAt
	} else {
		deletedAtTime = time.Time{}
	}
	return Table{
		ID:          entity.ID.ID,
		TableNumber: entity.TableNumber,
		CreatedAt:   entity.Timestamp.CreatedAt,
		UpdatedAt:   entity.Timestamp.UpdatedAt,
		DeletedAt: gorm.DeletedAt{
			Time:  deletedAtTime,
			Valid: entity.DeletedAt != nil,
		},
	}
}

func TableSchemaToEntity(schema Table) table.Table {
	return table.Table{
		ID:          identity.NewIDFromSchema(schema.ID),
		TableNumber: schema.TableNumber,
		Timestamp: shared.Timestamp{
			CreatedAt: schema.CreatedAt,
			UpdatedAt: schema.UpdatedAt,
			DeletedAt: &schema.DeletedAt.Time,
		},
	}
}
