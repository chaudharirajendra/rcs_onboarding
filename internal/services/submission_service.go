package services

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"rcs-onboarding/internal/models"
	"rcs-onboarding/internal/repositories"
	"rcs-onboarding/internal/utils"
)

type SubmissionService struct {
	subRepo   *repositories.SubmissionRepo
	formRepo  *repositories.FormRepo
	auditRepo *repositories.AuditRepo
}

func NewSubmissionService(subRepo *repositories.SubmissionRepo, formRepo *repositories.FormRepo, auditRepo *repositories.AuditRepo) *SubmissionService {
	return &SubmissionService{subRepo: subRepo, formRepo: formRepo, auditRepo: auditRepo}
}

func (s *SubmissionService) Submit(formType models.FormType, userID uint, dataStr string, isDraft bool) (*models.Submission, error) {
	template, err := s.formRepo.GetLatest(formType)
	if err != nil {
		return nil, err
	}

	validatedData, err := utils.ValidateData(template.Schema, dataStr, formType)
	if err != nil {
		return nil, err
	}

	status := models.Draft
	if !isDraft {
		status = models.Submitted
	}

	sub := &models.Submission{
		FormType:  formType,
		Version:   template.Version,
		UserID:    userID,
		Data:      validatedData,
		Status:    status,
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	if err := s.subRepo.Create(sub); err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *SubmissionService) UpdateDraft(id uint, userID uint, dataStr string) (*models.Submission, error) {
	sub, err := s.subRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if sub.UserID != userID || sub.Status != models.Draft {
		return nil, errors.New("unauthorized or invalid status")
	}

	template, err := s.formRepo.GetLatest(sub.FormType)
	if err != nil {
		return nil, err
	}

	validatedData, err := utils.ValidateData(template.Schema, dataStr, sub.FormType)
	if err != nil {
		return nil, err
	}

	sub.Data = validatedData
	sub.UpdatedBy = userID
	if err := s.subRepo.Update(sub); err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *SubmissionService) Review(id uint, userID uint, newStatus models.Status, remarks string) error {
	sub, err := s.subRepo.FindByID(id)
	if err != nil {
		return err
	}

	validTransitions := map[models.Status][]models.Status{
		models.Submitted: {models.InReview, models.Approved, models.Rejected},
		models.InReview:  {models.Approved, models.Rejected},
	}
	allowed, ok := validTransitions[sub.Status]
	if !ok || !containsStatus(allowed, newStatus) {
		return errors.New("invalid status transition")
	}

	sub.Status = newStatus
	sub.UpdatedBy = userID
	if err := s.subRepo.Update(sub); err != nil {
		return err
	}

	audit := &models.AuditLog{
		SubmissionID: sub.ID,
		UserID:       userID,
		Action:       fmt.Sprintf("Status changed to %s", newStatus),
		Remarks:      remarks,
	}
	return s.auditRepo.Create(audit)
}

func containsStatus(statuses []models.Status, target models.Status) bool {
	for _, st := range statuses {
		if st == target {
			return true
		}
	}
	return false
}

func (s *SubmissionService) GetFiltered(userID uint, role models.Role, customerID *uint, status *string, startDateStr *string, endDateStr *string, limitStr string, offsetStr string) ([]models.Submission, error) {
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(offsetStr)

	var startDate, endDate *time.Time
	if startDateStr != nil {
		t, _ := time.Parse("2006-01-02", *startDateStr)
		startDate = &t
	}
	if endDateStr != nil {
		t, _ := time.Parse("2006-01-02", *endDateStr)
		endDate = &t
	}

	return s.subRepo.FindFiltered(userID, role, customerID, status, startDate, endDate, limit, offset)
}

func (s *SubmissionService) GetByID(id uint, userID uint, role models.Role) (*models.Submission, error) {
	sub, err := s.subRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if role != models.Admin && sub.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	return sub, nil
}
