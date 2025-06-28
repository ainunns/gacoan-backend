package schema

import (
	"fp-kpl/domain/identity"
	"fp-kpl/domain/menu/category"
	"fp-kpl/domain/shared"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4();column:id"`
	Name      string         `gorm:"type:varchar(255);unique;not null;column:name"`
	CreatedAt time.Time      `gorm:"type:timestamp with time zone;column:created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp with time zone;column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp with time zone;column:deleted_at"`
}

func CategoryEntityToSchema(entity category.Category) Category {
	var deletedAtTime time.Time
	if entity.DeletedAt != nil {
		deletedAtTime = *entity.DeletedAt
	} else {
		deletedAtTime = time.Time{}
	}
	return Category{
		ID:        entity.ID.ID,
		Name:      entity.Name,
		CreatedAt: entity.Timestamp.CreatedAt,
		UpdatedAt: entity.Timestamp.UpdatedAt,
		DeletedAt: gorm.DeletedAt{
			Time:  deletedAtTime,
			Valid: entity.DeletedAt != nil,
		},
	}
}

func CategorySchemaToEntity(schema Category) category.Category {
	return category.Category{
		ID:   identity.NewIDFromSchema(schema.ID),
		Name: schema.Name,
		Timestamp: shared.Timestamp{
			CreatedAt: schema.CreatedAt,
			UpdatedAt: schema.UpdatedAt,
			DeletedAt: &schema.DeletedAt.Time,
		},
	}
}
