package repositories

import (
	"rcs-onboarding/internal/models"
	"time"

	"gorm.io/gorm"
)

type SubmissionRepo struct {
	db *gorm.DB
}

func NewSubmissionRepo(db *gorm.DB) *SubmissionRepo {
	return &SubmissionRepo{db: db}
}

func (r *SubmissionRepo) Create(sub *models.Submission) error {
	return r.db.Create(sub).Error
}

func (r *SubmissionRepo) FindByID(id uint) (*models.Submission, error) {
	var sub models.Submission
	err := r.db.First(&sub, id).Error
	return &sub, err
}

func (r *SubmissionRepo) Update(sub *models.Submission) error {
	return r.db.Save(sub).Error
}

func (r *SubmissionRepo) FindFiltered(userID uint, role models.Role, customerID *uint, status *string, startDate *time.Time, endDate *time.Time, limit int, offset int) ([]models.Submission, error) {
	query := r.db.Limit(limit).Offset(offset)
	if role != models.Admin {
		query = query.Where("user_id = ?", userID)
	}
	if customerID != nil {
		query = query.Where("user_id = ?", *customerID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if startDate != nil && endDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", *startDate, *endDate)
	}
	var subs []models.Submission
	err := query.Find(&subs).Error
	return subs, err
}
