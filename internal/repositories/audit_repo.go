package repositories

import (
	"rcs-onboarding/internal/models"

	"gorm.io/gorm"
)

type AuditRepo struct {
	db *gorm.DB
}

func NewAuditRepo(db *gorm.DB) *AuditRepo {
	return &AuditRepo{db: db}
}

func (r *AuditRepo) Create(audit *models.AuditLog) error {
	return r.db.Create(audit).Error
}
