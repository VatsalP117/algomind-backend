package handlers

import (
	"net/http"

	"github.com/VatsalP117/algomind-backend/internal/database"
	"github.com/VatsalP117/algomind-backend/internal/models"
	"github.com/labstack/echo/v4"
)

type ConceptHandler struct {
	DB *database.Service
}

func NewConceptHandler(db *database.Service) *ConceptHandler {
	return &ConceptHandler{DB: db}
}

// ListConcepts returns Global Concepts + The User's Custom Concepts
func (h *ConceptHandler) ListConcepts(c echo.Context) error {
	// 1. Define the target slice
	var concepts []models.Concept

	// 2. The Query
	query := `SELECT * FROM concepts WHERE user_id IS NULL ORDER BY title ASC`

	// 3. One-Liner Selection (The "Magic")
	// .Select() automatically iterates rows and maps columns to struct tags
	err := h.DB.Db.SelectContext(c.Request().Context(), &concepts, query)
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// Safety: If no rows found, sqlx returns empty slice (not nil), which is JSON []
	if concepts == nil {
		concepts = []models.Concept{}
	}

	return c.JSON(http.StatusOK, concepts)
}