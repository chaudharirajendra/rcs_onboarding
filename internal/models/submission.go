package models

import "gorm.io/gorm"

type Submission struct {
	gorm.Model
	FormType  FormType
	Version   int
	UserID    uint
	Data      string `gorm:"type:text"` // JSON map[string]any
	Status    Status
	CreatedBy uint
	UpdatedBy uint
}
