package models

import (
	"time"
)

// Concept: The "Theory" (BFS, Sliding Window, etc.)
type Concept struct {
    // Add the `db` tag matching your SQL column names
    ID          int64     `json:"id"          db:"id"`
    UserID      *int64    `json:"user_id"     db:"user_id"`
    Title       string    `json:"title"       db:"title"`
    Description string    `json:"description" db:"description"`
    Content     string    `json:"content"     db:"content"`
    CreatedAt   time.Time `json:"created_at"  db:"created_at"`
}

// UserItem: The tracking row for SRS (Links to Concept or is a Problem)
type UserItem struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	ItemType     string    `json:"item_type"` // "CONCEPT" or "PROBLEM"
	ConceptID    *int64    `json:"concept_id,omitempty"`
	ProblemTitle string    `json:"problem_title,omitempty"`
	ProblemLink  string    `json:"problem_link,omitempty"`
	
	// SRS State
	NextReviewAt time.Time `json:"next_review_at"`
	IntervalDays int       `json:"interval_days"`
	EaseFactor   float64   `json:"ease_factor"`
	Streak       int       `json:"streak"`
}

// ReviewLog: A history entry for the heatmap
type ReviewLog struct {
	ID         int64     `json:"id"`
	UserItemID int64     `json:"user_item_id"`
	Rating     string    `json:"rating"` // AGAIN, HARD, GOOD, EASY
	ReviewedAt time.Time `json:"reviewed_at"`
}

// Incoming JSON for creating a new item
type CreateItemRequest struct {
	Type         string `json:"type" validate:"required,oneof=PROBLEM CONCEPT"` // Validator ensures it's one of these two
	ConceptID    *int64 `json:"concept_id"` // Nullable (if it's a root concept)
	ProblemTitle string `json:"problem_title"`
	ProblemLink  string `json:"problem_link"`
}