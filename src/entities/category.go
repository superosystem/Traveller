package entities

import (
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"time"
)

type Category struct {
	ID        string    `gorm:"primaryKey;size:200" json:"id"`
	Name      string    `gorm:"size:255" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (rec *Category) ToCategoryDomain() *domain.Category {
	return &domain.Category{
		ID:        rec.ID,
		Name:      rec.Name,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
	}
}

func FromCategoryDomain(categoryDomain *domain.Category) *Category {
	return &Category{
		ID:        categoryDomain.ID,
		Name:      categoryDomain.Name,
		CreatedAt: categoryDomain.CreatedAt,
		UpdatedAt: categoryDomain.UpdatedAt,
	}
}
