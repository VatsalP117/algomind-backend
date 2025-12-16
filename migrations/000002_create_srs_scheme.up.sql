-- 1. CONCEPTS: Theory definitions (Global OR User-specific)
CREATE TABLE concepts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE, -- NULL = System/Global, Value = Custom
    title TEXT NOT NULL,
    description TEXT,     -- Short summary
    content TEXT,         -- Full Markdown explanation
    created_at TIMESTAMP(0) WITH TIME ZONE DEFAULT NOW()
);

-- Index to quickly find "My Concepts" OR "System Concepts"
CREATE INDEX idx_concepts_user ON concepts(user_id);


-- 2. USER_ITEMS: The central SRS tracking table
-- This stores BOTH "Problems" (LeetCode) AND "Concepts" (Theory)
CREATE TABLE user_items (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- What kind of item is this?
    item_type TEXT NOT NULL CHECK (item_type IN ('CONCEPT', 'PROBLEM')),
    
    -- Links (Polymorphic-ish relationship)
    concept_id BIGINT REFERENCES concepts(id), -- If it's a CONCEPT, this links to the definition. If PROBLEM, links to parent.
    
    -- Problem Specifics (Only used if item_type = 'PROBLEM')
    problem_title TEXT,
    problem_link TEXT,
    
    -- SRS Algorithm Data (Spaced Repetition State)
    next_review_at TIMESTAMP(0) WITH TIME ZONE DEFAULT NOW(),
    interval_days INTEGER DEFAULT 0, -- How many days until next review
    ease_factor NUMERIC(4, 2) DEFAULT 2.50, -- Standard SM-2 starting ease
    streak INTEGER DEFAULT 0,
    
    created_at TIMESTAMP(0) WITH TIME ZONE DEFAULT NOW()
);

-- Index for the "Daily Review Queue" (Find items due today for this user)
CREATE INDEX idx_reviews_due ON user_items(user_id, next_review_at);
-- Index to find all problems belonging to a specific concept (for the Cascading Reset)
CREATE INDEX idx_items_concept ON user_items(user_id, concept_id);


-- 3. REVIEW_LOGS: History for Heatmaps/Analytics
CREATE TABLE review_logs (
    id BIGSERIAL PRIMARY KEY,
    user_item_id BIGINT NOT NULL REFERENCES user_items(id) ON DELETE CASCADE,
    rating TEXT NOT NULL CHECK (rating IN ('AGAIN', 'HARD', 'GOOD', 'EASY')),
    reviewed_at TIMESTAMP(0) WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_logs_date ON review_logs(reviewed_at);