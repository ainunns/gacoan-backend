package schema

import (
	"database/sql/driver"
	"fmt"
	"fp-kpl/domain/identity"
	menu "fp-kpl/domain/menu/menu_item"
	"fp-kpl/domain/shared"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Duration struct {
	time.Duration
}

func (d *Duration) Scan(value interface{}) error {
	if value == nil {
		d.Duration = 0
		return nil
	}

	switch v := value.(type) {
	case string:
		return d.parseInterval(v)
	case []byte:
		return d.parseInterval(string(v))
	case int64:
		d.Duration = time.Duration(v)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into Duration", value)
	}
}

func (d Duration) Value() (driver.Value, error) {
	// Convert to PostgreSQL interval format instead of nanoseconds
	// Format: HH:MM:SS or HH:MM:SS.microseconds
	hours := int(d.Duration.Hours())
	minutes := int(d.Duration.Minutes()) % 60
	seconds := int(d.Duration.Seconds()) % 60
	microseconds := int(d.Duration.Microseconds()) % 1000000

	if microseconds > 0 {
		return fmt.Sprintf("%02d:%02d:%02d.%06d", hours, minutes, seconds, microseconds), nil
	}
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds), nil
}

func (d *Duration) parseInterval(interval string) error {
	interval = strings.TrimSpace(interval)

	if interval == "" {
		d.Duration = 0
		return nil
	}

	parts := strings.Split(interval, ":")
	if len(parts) != 3 {
		return fmt.Errorf("invalid interval format: %s", interval)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid hours in interval: %s", parts[0])
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid minutes in interval: %s", parts[1])
	}

	secondsStr := parts[2]
	var seconds int
	var microseconds int

	if strings.Contains(secondsStr, ".") {
		secondsParts := strings.Split(secondsStr, ".")
		if len(secondsParts) != 2 {
			return fmt.Errorf("invalid seconds format: %s", secondsStr)
		}

		seconds, err = strconv.Atoi(secondsParts[0])
		if err != nil {
			return fmt.Errorf("invalid seconds in interval: %s", secondsParts[0])
		}

		microStr := secondsParts[1]
		if len(microStr) > 6 {
			microStr = microStr[:6]
		} else if len(microStr) < 6 {
			microStr = microStr + strings.Repeat("0", 6-len(microStr))
		}

		microseconds, err = strconv.Atoi(microStr)
		if err != nil {
			return fmt.Errorf("invalid microseconds in interval: %s", secondsParts[1])
		}
	} else {
		seconds, err = strconv.Atoi(secondsStr)
		if err != nil {
			return fmt.Errorf("invalid seconds in interval: %s", secondsStr)
		}
	}

	d.Duration = time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second +
		time.Duration(microseconds)*time.Microsecond

	return nil
}

type Menu struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();column:id"`
	CategoryID  uuid.UUID       `gorm:"type:uuid;not null;column:category_id"`
	Name        string          `gorm:"type:varchar(255);not null;column:name"`
	ImageURL    string          `gorm:"type:varchar(255);not null;column:image_url"`
	Price       decimal.Decimal `gorm:"type:decimal(10,2);not null;column:price"`
	IsAvailable bool            `gorm:"type:boolean;not null;column:is_available"`
	CookingTime Duration        `gorm:"type:interval;not null;column:cooking_time"`
	Description string          `gorm:"type:text;not null;column:description"`
	CreatedAt   time.Time       `gorm:"type:timestamp with time zone;not null;column:created_at"`
	UpdatedAt   time.Time       `gorm:"type:timestamp with time zone;not null;column:updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"type:timestamp with time zone;column:deleted_at"`

	Category *Category `gorm:"foreignKey:CategoryID"`
	Orders   []Order   `gorm:"foreignKey:MenuID"`
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
		CookingTime: Duration{Duration: entity.CookingTime},
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
	return menu.Menu{
		ID:          identity.NewIDFromSchema(schema.ID),
		CategoryID:  identity.NewIDFromSchema(schema.CategoryID),
		Name:        schema.Name,
		ImageURL:    shared.NewURLFromSchema(schema.ImageURL),
		Price:       shared.NewPriceFromSchema(schema.Price),
		IsAvailable: schema.IsAvailable,
		CookingTime: schema.CookingTime.Duration,
		Description: schema.Description,
		Timestamp: shared.Timestamp{
			CreatedAt: schema.CreatedAt,
			UpdatedAt: schema.UpdatedAt,
			DeletedAt: &schema.DeletedAt.Time,
		},
	}
}
