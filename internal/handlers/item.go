package handlers

import (
	"net/http"
	"time"

	"github.com/VatsalP117/algomind-backend/internal/database"
	"github.com/VatsalP117/algomind-backend/internal/models"
	"github.com/labstack/echo/v4"
)

type ItemHandler struct {
	DB *database.Service
}

func NewItemHandler(db *database.Service) *ItemHandler {
	return &ItemHandler{DB: db}
}

func (h *ItemHandler) CreateItem(c echo.Context) error {
	// 1. Parse the Request Body
	var req models.CreateItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	// 2. Get the Internal User ID
	// (We need a helper for this because our tables assume Int64, not Clerk Strings)
	clerkID := c.Get("user_id").(string)
	userID, err := h.getInternalUserID(c, clerkID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not found"})
	}

	// 3. Prepare the Insert
	// Notice we set defaults: Next Review is NOW, Interval is 0 days.
	query := `
		INSERT INTO user_items (
			user_id, item_type, concept_id, problem_title, problem_link, 
			next_review_at, interval_days, ease_factor, streak
		) VALUES (
			:user_id, :item_type, :concept_id, :problem_title, :problem_link,
			:next_review_at, 0, 2.5, 0
		) RETURNING id`

	// Construct the data map for sqlx.NamedQuery
	// We use a map or struct. Here, a quick anonymous struct is clean.
	params := map[string]interface{}{
		"user_id":        userID,
		"item_type":      req.Type,
		"concept_id":     req.ConceptID,
		"problem_title":  req.ProblemTitle,
		"problem_link":   req.ProblemLink,
		"next_review_at": time.Now(), // Due immediately!
	}

	// 4. Execute
	// NamedQuery is an sqlx feature that lets you use :name instead of $1, $2
	// It's much safer and easier to read.
	rows, err := h.DB.Db.NamedQueryContext(c.Request().Context(), query, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save item: " + err.Error()})
	}
	defer rows.Close()

	// Get the generated ID
	var newID int64
	if rows.Next() {
		rows.Scan(&newID)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":      newID,
		"message": "Item added to review queue",
	})
}

// Helper: Resolves Clerk String ID -> DB Int64 ID
func (h *ItemHandler) getInternalUserID(c echo.Context, clerkID string) (int64, error) {
	var id int64
	query := `SELECT id FROM users WHERE clerk_id = $1`
	err := h.DB.Db.GetContext(c.Request().Context(), &id, query, clerkID)
	return id, err
}