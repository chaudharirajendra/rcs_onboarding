package models

import "gorm.io/gorm"

type FormVersion struct {
	gorm.Model
	Type    FormType `gorm:"index"`
	Version int      `gorm:"index"`
	Schema  string   `gorm:"type:text"` // JSON []Field
}

type Field struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"` // string, int, url, email, lookup
	Required bool     `json:"required"`
	Max      int      `json:"max,omitempty"`
	Min      int      `json:"min,omitempty"`
	Options  []string `json:"options,omitempty"`
}
