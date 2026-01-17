package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string // Hashed
	Role     Role
}

type FormType string

const (
	CustomerOrder FormType = "customer_order"
	Qualification FormType = "qualification"
)

type Status string

const (
	Draft     Status = "Draft"
	Submitted Status = "Submitted"
	InReview  Status = "In Review"
	Approved  Status = "Approved"
	Rejected  Status = "Rejected"
)

type Role string

const (
	Customer Role = "customer"
	TPM      Role = "tpm"
	Sales    Role = "sales"
	Admin    Role = "admin"
)
