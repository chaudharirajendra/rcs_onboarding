package models

import "gorm.io/gorm"

type AuditLog struct {
	gorm.Model
	SubmissionID uint
	UserID       uint
	Action       string
	Remarks      string
}
