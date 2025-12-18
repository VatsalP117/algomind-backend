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


type UserItem struct {
	ID           int64     `json:"id"             db:"id"`
	UserID       int64     `json:"user_id"        db:"user_id"`
	ItemType     string    `json:"item_type"      db:"item_type"`
	ConceptID    *int64    `json:"concept_id"     db:"concept_id"`
	ProblemTitle string    `json:"problem_title"  db:"problem_title"`
	ProblemLink  string    `json:"problem_link"   db:"problem_link"`
	
	// SRS State
	NextReviewAt time.Time `json:"next_review_at" db:"next_review_at"`
	IntervalDays int       `json:"interval_days"  db:"interval_days"`
	EaseFactor   float64   `json:"ease_factor"    db:"ease_factor"`
	Streak       int       `json:"streak"         db:"streak"`
	
	// Crucial: SELECT * returns created_at, so we must map it!
	CreatedAt    time.Time `json:"created_at"     db:"created_at"` 
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

// ReviewQueueItem: A combined view for the frontend "Flashcard"
type ReviewQueueItem struct {
	UserItem // Embeds all the tracking fields (ID, NextReviewAt, etc.)

	// Extra fields we get by JOINing with the concepts table
	ConceptTitle   *string `db:"concept_title"   json:"concept_title,omitempty"`
	ConceptContent *string `db:"concept_content" json:"concept_content,omitempty"`
}