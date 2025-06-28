package schema

import (
	"fp-kpl/domain/identity"
	menu "fp-kpl/domain/menu/menu_item"
	"fp-kpl/domain/shared"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Menu struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();column:id"`
	CategoryID  uuid.UUID       `gorm:"type:uuid;not null;column:category_id"`
	Name        string          `gorm:"type:varchar(255);not null;column:name"`
	ImageURL    string          `gorm:"type:varchar(255);not null;column:image_url"`
	Price       decimal.Decimal `gorm:"type:decimal(10,2);not null;column:price"`
	IsAvailable bool            `gorm:"type:boolean;not null;column:is_available"`
	CookingTime time.Time       `gorm:"type:timestamp with time zone;not null;column:cooking_time"`
	Description string          `gorm:"type:text;not null;column:description"`
	CreatedAt   time.Time       `gorm:"type:timestamp with time zone;not null;column:created_at"`
	UpdatedAt   time.Time       `gorm:"type:timestamp with time zone;not null;column:updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"type:timestamp with time zone;column:deleted_at"`

	Category *Category `gorm:"foreignKey:CategoryID"`
}

func MenuEntityToSchema(entity menu.Menu) Menu {
	var deletedAtTime time.Time
	if entity.DeletedAt != nil {
		deletedAtTime = *entity.DeletedAt
	} else {
		deletedAtTime = time.Time{}
	}

	return Menu{
		ID:          entity.ID.ID,
		Name:        entity.Name,
		ImageURL:    entity.ImageURL.Path,
		Price:       entity.Price.Price,
		IsAvailable: entity.IsAvailable,
		CookingTime: entity.CookingTime,
		Description: entity.Description,
		CreatedAt:   entity.Timestamp.CreatedAt,
		UpdatedAt:   entity.Timestamp.UpdatedAt,
		DeletedAt: gorm.DeletedAt{
			Time:  deletedAtTime,
			Valid: entity.DeletedAt != nil,
		},
		CategoryID: entity.CategoryID.ID,
	}
}

func MenuSchemaToEntity(schema Menu) menu.Menu {
	var deletedAtTime time.Time
	if schema.DeletedAt.Valid {
		deletedAtTime = schema.DeletedAt.Time
	} else {
		deletedAtTime = time.Time{}
	}

	return menu.Menu{
		ID:          identity.NewIDFromSchema(schema.ID),
		CategoryID:  identity.NewIDFromSchema(schema.CategoryID),
		Name:        schema.Name,
		ImageURL:    shared.NewURLFromSchema(schema.ImageURL),
		Price:       shared.NewPriceFromSchema(schema.Price),
		IsAvailable: schema.IsAvailable,
		CookingTime: schema.CookingTime,
		Description: schema.Description,
		Timestamp: shared.Timestamp{
			CreatedAt: schema.CreatedAt,
			UpdatedAt: schema.UpdatedAt,
			DeletedAt: &deletedAtTime,
		},
	}
}
