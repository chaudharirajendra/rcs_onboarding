package repositories

import (
	"rcs-onboarding/internal/models"

	"gorm.io/gorm"
)

type FormRepo struct {
	db *gorm.DB
}

func NewFormRepo(db *gorm.DB) *FormRepo {
	return &FormRepo{db: db}
}

func (r *FormRepo) Create(form *models.FormVersion) error {
	return r.db.Create(form).Error
}

func (r *FormRepo) GetLatest(formType models.FormType) (*models.FormVersion, error) {
	var template models.FormVersion
	err := r.db.Where("type = ?", formType).Order("version desc").First(&template).Error
	return &template, err
}

func (r *FormRepo) ListVersions(formType models.FormType) ([]models.FormVersion, error) {
	var templates []models.FormVersion
	err := r.db.Where("type = ?", formType).Order("version desc").Find(&templates).Error
	return templates, err
}
