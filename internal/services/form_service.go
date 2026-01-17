package services

import (
	"encoding/json"
	"errors"
	"rcs-onboarding/internal/models"
	"rcs-onboarding/internal/repositories"

	"gorm.io/gorm"
)

type FormService struct {
	repo *repositories.FormRepo
}

func NewFormService(repo *repositories.FormRepo) *FormService {
	return &FormService{repo: repo}
}

func (s *FormService) Create(formType models.FormType, schema []models.Field) (*models.FormVersion, error) {
	// Compute next version
	latest, err := s.repo.GetLatest(formType)
	nextVersion := 1
	if err == nil {
		nextVersion = latest.Version + 1
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	schemaJSON, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	newForm := &models.FormVersion{
		Type:    formType,
		Version: nextVersion,
		Schema:  string(schemaJSON),
	}

	if err := s.repo.Create(newForm); err != nil {
		return nil, err
	}

	return newForm, nil
}

func (s *FormService) GetLatest(formType models.FormType) (*models.FormVersion, error) {
	return s.repo.GetLatest(formType)
}

func (s *FormService) ListVersions(formType models.FormType) ([]models.FormVersion, error) {
	return s.repo.ListVersions(formType)
}
