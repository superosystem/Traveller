package materials

import (
	"time"

	"github.com/superosystem/TrainingSystem/backend/domain/materials"
	"github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/modules"
	"gorm.io/gorm"
)

type Material struct {
	ID          string `json:"id" gorm:"primaryKey;size:200"`
	ModuleId    string `json:"module_id" gorm:"size:200"`
	Title       string `json:"title" gorm:"size:255"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Module      modules.Module
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (rec *Material) ToDomain() *materials.Domain {
	return &materials.Domain{
		ID:          rec.ID,
		CourseId:    rec.Module.CourseId,
		ModuleId:    rec.ModuleId,
		Title:       rec.Title,
		URL:         rec.URL,
		Description: rec.Description,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
		DeletedAt:   rec.DeletedAt,
	}
}

func FromDomain(materialDomain *materials.Domain) *Material {
	return &Material{
		ID:          materialDomain.ID,
		ModuleId:    materialDomain.ModuleId,
		Title:       materialDomain.Title,
		URL:         materialDomain.URL,
		Description: materialDomain.Description,
		CreatedAt:   materialDomain.CreatedAt,
		UpdatedAt:   materialDomain.UpdatedAt,
		DeletedAt:   materialDomain.DeletedAt,
	}
}
