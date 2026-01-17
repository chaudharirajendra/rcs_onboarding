package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"rcs-onboarding/internal/models"
	"rcs-onboarding/internal/services"

	"github.com/gin-gonic/gin"
)

type SubmissionHandler struct {
	subService   *services.SubmissionService
	auditService *services.AuditService
}

func NewSubmissionHandler(subService *services.SubmissionService, auditService *services.AuditService) *SubmissionHandler {
	return &SubmissionHandler{subService: subService, auditService: auditService}
}

// ... (imports remain the same)

func (h *SubmissionHandler) Submit(c *gin.Context) {
	formTypeStr := c.Param("id") // Now using "id" as the param key
	var formType models.FormType
	switch formTypeStr {
	case "customer_order":
		formType = models.CustomerOrder
	case "qualification":
		formType = models.Qualification
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form type"})
		return
	}

	userID := c.GetUint("userID")

	var req struct {
		Data    json.RawMessage `json:"data"`
		IsDraft bool            `json:"is_draft"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.subService.Submit(formType, userID, string(req.Data), req.IsDraft)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !req.IsDraft {
		h.auditService.CreateAudit(sub.ID, userID, "Submitted", "Initial submission")
	}

	c.JSON(http.StatusCreated, sub)
}

func (h *SubmissionHandler) Review(c *gin.Context) {
	idStr := c.Param("id") // Now using "id" as the param key
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid submission ID (must be numeric)"})
		return
	}
	userID := c.GetUint("userID")

	var req struct {
		Status  models.Status `json:"status"`
		Remarks string        `json:"remarks"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.subService.Review(uint(id), userID, req.Status, req.Remarks)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "review completed"})
}

func (h *SubmissionHandler) UpdateDraft(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	userID := c.GetUint("userID")

	var req struct {
		Data json.RawMessage `json:"data"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.subService.UpdateDraft(uint(id), userID, string(req.Data))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.auditService.CreateAudit(sub.ID, userID, "Updated Draft", "Draft updated")

	c.JSON(http.StatusOK, sub)
}

func (h *SubmissionHandler) GetFiltered(c *gin.Context) {
	userID := c.GetUint("userID")
	role := models.Role(c.GetString("role"))

	var customerID *uint
	if cidStr := c.Query("customer_id"); cidStr != "" && role == models.Admin {
		cid64, err := strconv.ParseUint(cidStr, 10, 32)
		if err == nil {
			cid := uint(cid64)
			customerID = &cid
		}
	}

	status := getStringPtr(c.Query("status"))
	startDate := getStringPtr(c.Query("start_date"))
	endDate := getStringPtr(c.Query("end_date"))
	limit := c.Query("limit")
	offset := c.Query("offset")

	subs, err := h.subService.GetFiltered(userID, role, customerID, status, startDate, endDate, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subs)
}

func (h *SubmissionHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	userID := c.GetUint("userID")
	role := models.Role(c.GetString("role"))

	sub, err := h.subService.GetByID(uint(id), userID, role)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sub)
}

func getStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
