package handlers

import (
	"database/sql"
	"net/http"

	"github.com/VatsalP117/algomind-backend/internal/database"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	DB *database.Service
}

func NewUserHandler(db *database.Service) *UserHandler {
	return &UserHandler{DB: db}
}

// GetProfile finds the user in the DB. If they don't exist, it creates them.
func (h *UserHandler) GetProfile(c echo.Context) error {
	// 1. Get the authenticated Clerk ID from the context
	clerkID := c.Get("user_id").(string)

	// 2. Try to find the user in OUR database
	// We use a struct to hold the result (using sqlx tags if you added them, or manual scanning)
	var user struct {
		ID      int64  `db:"id"`
		ClerkID string `db:"clerk_id"`
		Email   string `db:"email"`
	}

	// 'Get' is a helper from sqlx that selects a single row
	query := `SELECT id, clerk_id, email FROM users WHERE clerk_id = $1`
	err := h.DB.Db.Get(&user, query, clerkID)

	// 3. Handle the result
	if err == sql.ErrNoRows {
		// --- CASE A: User does not exist (First Login) ---
		// We insert them now.
		// Note: In a real app, you might want to fetch the real email from Clerk's API.
		// For now, we use a placeholder so the constraint doesn't fail.
		defaultEmail := "user_" + clerkID + "@placeholder.com"

		insertQuery := `INSERT INTO users (clerk_id, email) VALUES ($1, $2) RETURNING id, clerk_id, email`
		
		// We use QueryRow because we want the returned ID
		err = h.DB.Db.QueryRowx(insertQuery, clerkID, defaultEmail).StructScan(&user)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to register user: " + err.Error()})
		}
		
	} else if err != nil {
		// --- CASE B: Real Database Error ---
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// 4. Return the Profile (Case C: User existed, or was just created)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "Profile fetched successfully",
		"internal_id": user.ID,      // This is what we need for foreign keys!
		"clerk_id":    user.ClerkID,
		"email":       user.Email,
	})
}