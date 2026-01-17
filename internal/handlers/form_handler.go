package handlers

import (
	"net/http"

	"rcs-onboarding/internal/models"
	"rcs-onboarding/internal/services"

	"github.com/gin-gonic/gin"
)

type FormHandler struct {
	service *services.FormService
}

func NewFormHandler(service *services.FormService) *FormHandler {
	return &FormHandler{service: service}
}

func (h *FormHandler) Create(c *gin.Context) {
	formType := models.FormType(c.Param("type"))
	var req struct {
		Schema []models.Field `json:"schema" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schema: " + err.Error()})
		return
	}

	// Basic validation (extend as needed)
	if len(req.Schema) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Schema cannot be empty"})
		return
	}

	newVersion, err := h.service.Create(formType, req.Schema)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newVersion)
}

func (h *FormHandler) ListVersions(c *gin.Context) {
	formType := models.FormType(c.Param("type"))
	templates, err := h.service.ListVersions(formType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, templates)
}

func (h *FormHandler) GetLatest(c *gin.Context) {
	formType := models.FormType(c.Param("type"))
	template, err := h.service.GetLatest(formType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, template)
}
