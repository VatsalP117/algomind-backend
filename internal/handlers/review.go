package handlers

import (
	"net/http"

	"github.com/VatsalP117/algomind-backend/internal/database"
	"github.com/VatsalP117/algomind-backend/internal/models"
	"github.com/labstack/echo/v4"
)

type ReviewHandler struct {
	DB *database.Service
}

func NewReviewHandler(db *database.Service) *ReviewHandler {
	return &ReviewHandler{DB: db}
}

// GetQueue returns all items due for review (NextReviewAt <= NOW)
func (h *ReviewHandler) GetQueue(c echo.Context) error {
	clerkID := c.Get("user_id").(string)

	// We need the internal ID for the query
	var userID int64
	err := h.DB.Db.GetContext(c.Request().Context(), &userID, "SELECT id FROM users WHERE clerk_id=$1", clerkID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not found"})
	}

	// The Query:
	// 1. Selects everything from user_items
	// 2. Joins 'concepts' to get the Title/Content (Aliased as concept_*)
	// 3. Filters by User AND Due Date
	query := `
		SELECT 
			ui.*,
			c.title AS concept_title,
			c.content AS concept_content
		FROM user_items ui
		LEFT JOIN concepts c ON ui.concept_id = c.id
		WHERE ui.user_id = $1 
		  AND ui.next_review_at <= NOW()
		ORDER BY ui.next_review_at ASC
		LIMIT 50
	`

	var queue []models.ReviewQueueItem
	err = h.DB.Db.SelectContext(c.Request().Context(), &queue, query, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch queue: " + err.Error()})
	}

	// Handle empty result
	if queue == nil {
		queue = []models.ReviewQueueItem{}
	}

	return c.JSON(http.StatusOK, queue)
}