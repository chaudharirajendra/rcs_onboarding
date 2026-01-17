package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"rcs-onboarding/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func ValidateData(schemaStr string, dataStr string, formType models.FormType) (string, error) {
	var schema []models.Field
	if err := json.Unmarshal([]byte(schemaStr), &schema); err != nil {
		return "", err
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		return "", err
	}

	for _, f := range schema {
		val, ok := data[f.Name]
		if !ok && f.Required {
			return "", fmt.Errorf("%s is required", f.Name)
		}
		if !ok {
			continue
		}

		switch f.Type {
		case "string":
			s, ok := val.(string)
			if !ok {
				return "", fmt.Errorf("%s must be string", f.Name)
			}
			if f.Max > 0 && len(s) > f.Max {
				return "", fmt.Errorf("%s exceeds max length %d", f.Name, f.Max)
			}
			if f.Min > 0 && len(s) < f.Min {
				return "", fmt.Errorf("%s below min length %d", f.Name, f.Min)
			}
		case "int":
			iStr, ok := val.(string)
			if !ok {
				i, ok := val.(float64)
				if !ok {
					return "", fmt.Errorf("%s must be integer", f.Name)
				}
				iStr = strconv.Itoa(int(i))
			}
			i, err := strconv.Atoi(iStr)
			if err != nil {
				return "", fmt.Errorf("%s must be integer", f.Name)
			}
			if f.Min > 0 && i < f.Min {
				return "", fmt.Errorf("%s below min %d", f.Name, f.Min)
			}
			if f.Max > 0 && i > f.Max {
				return "", fmt.Errorf("%s exceeds max %d", f.Name, f.Max)
			}
		case "url":
			s, ok := val.(string)
			if !ok {
				return "", fmt.Errorf("%s must be string", f.Name)
			}
			_, err := url.ParseRequestURI(s)
			if err != nil {
				return "", fmt.Errorf("%s invalid URL", f.Name)
			}
		case "email":
			s, ok := val.(string)
			if !ok {
				return "", fmt.Errorf("%s must be string", f.Name)
			}
			if !emailRegex.MatchString(s) {
				return "", fmt.Errorf("%s invalid email", f.Name)
			}
			if f.Max > 0 && len(s) > f.Max {
				return "", fmt.Errorf("%s exceeds max length %d", f.Name, f.Max)
			}
		case "lookup":
			s, ok := val.(string)
			if !ok {
				return "", fmt.Errorf("%s must be string", f.Name)
			}
			if !contains(f.Options, s) {
				return "", fmt.Errorf("%s invalid option: %s", f.Name, s)
			}
		default:
			return "", fmt.Errorf("unknown type %s for %s", f.Type, f.Name)
		}
	}

	// Business rule example
	if zip, ok := data["address_zip_code"]; ok {
		if s, ok := zip.(string); ok && !strings.Contains(s, "0123456789") {
			return "", errors.New("address_zip_code must be numeric")
		}
	}

	// Auto SID
	if formType == models.CustomerOrder && data["sid"] == nil {
		data["sid"] = "HSN-" + strings.ToUpper(uuid.New().String()[:8])
	}

	updatedData, _ := json.Marshal(data)
	return string(updatedData), nil
}

func contains(options []string, val string) bool {
	for _, o := range options {
		if strings.EqualFold(o, val) {
			return true
		}
	}
	return false
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Seed functions placed here for structure compliance
func SeedTemplates(db *gorm.DB) {
	var count int64
	db.Model(&models.FormVersion{}).Count(&count)
	if count > 0 {
		return
	}

	// Qualification Schema
	qualFields := []models.Field{
		{Name: "organization_name", Type: "string", Required: true, Max: 100},
		{Name: "legal_company_name", Type: "string", Required: true, Max: 100},
		{Name: "business_ein", Type: "string", Required: true, Max: 20},
		{Name: "organization_type", Type: "lookup", Required: true, Options: []string{"public", "private", "non_profit", "government"}},
		{Name: "address_line_1", Type: "string", Required: true, Max: 100},
		{Name: "address_line_2", Type: "string", Required: false, Max: 100},
		{Name: "address_city", Type: "string", Required: true, Max: 100},
		{Name: "address_state", Type: "string", Required: true, Max: 2},
		{Name: "address_zip_code", Type: "string", Required: true, Min: 6, Max: 6},
		{Name: "phone_number", Type: "string", Required: true, Min: 10, Max: 12},
		{Name: "phone_label", Type: "string", Required: true, Max: 50},
		{Name: "website_url", Type: "url", Required: true},
		{Name: "website_label", Type: "string", Required: true, Max: 50},
		{Name: "terms_conditions", Type: "url", Required: true},
		{Name: "privacy_policy", Type: "url", Required: true},
		{Name: "contact_name", Type: "string", Required: true, Max: 100},
		{Name: "contact_position", Type: "string", Required: true, Max: 100},
		{Name: "contact_email", Type: "email", Required: true, Max: 100},
		{Name: "contact_phone_number", Type: "string", Required: true, Min: 10, Max: 12},
		{Name: "documents", Type: "string", Required: false, Max: 500},
	}
	qualSchema, _ := json.Marshal(qualFields)
	db.Create(&models.FormVersion{Type: models.Qualification, Version: 1, Schema: string(qualSchema)})

	// Customer Order Schema
	orderFields := []models.Field{
		{Name: "brand_name", Type: "string", Required: true, Max: 50},
		{Name: "brand_description", Type: "string", Required: true, Max: 200},
		{Name: "brand_tagline", Type: "string", Required: true, Max: 100},
		{Name: "brand_logo_image", Type: "url", Required: true},
		{Name: "banner_image", Type: "url", Required: true},
		{Name: "agent_name", Type: "string", Required: true, Max: 50},
		{Name: "agent_purpose", Type: "lookup", Required: true, Options: []string{"OTP", "TRANSACTION", "PROMOTION", "ALERTS", "CUSTOMER_SERVICE"}},
		{Name: "agent_billing_category", Type: "lookup", Required: true, Options: []string{"BASIC_MESSAGE", "PREMIUM_MESSAGE"}},
		{Name: "agent_service_code", Type: "string", Required: true, Max: 20},
		{Name: "color", Type: "string", Required: true, Max: 10},
		{Name: "languages", Type: "string", Required: true, Max: 100},
		{Name: "message_webhook_url", Type: "url", Required: true},
		{Name: "sid", Type: "string", Required: false},
	}
	orderSchema, _ := json.Marshal(orderFields)
	db.Create(&models.FormVersion{Type: models.CustomerOrder, Version: 1, Schema: string(orderSchema)})
}

func SeedUsers(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	db.Create(&models.User{Username: "admin", Password: string(hash), Role: models.Admin})
	db.Create(&models.User{Username: "customer", Password: string(hash), Role: models.Customer})
	db.Create(&models.User{Username: "tpm", Password: string(hash), Role: models.TPM})
	db.Create(&models.User{Username: "sales", Password: string(hash), Role: models.Sales})
}
