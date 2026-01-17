package services

import (
	"rcs-onboarding/internal/models"
	"rcs-onboarding/internal/repositories"
)

type AuditService struct {
	repo *repositories.AuditRepo
}

func NewAuditService(repo *repositories.AuditRepo) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) CreateAudit(submissionID uint, userID uint, action string, remarks string) error {
	audit := &models.AuditLog{
		SubmissionID: submissionID,
		UserID:       userID,
		Action:       action,
		Remarks:      remarks,
	}
	return s.repo.Create(audit)
}
