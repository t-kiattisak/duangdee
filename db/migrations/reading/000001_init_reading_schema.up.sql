-- Up Migration for reading_db
CREATE TABLE IF NOT EXISTS reading_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    spread_id VARCHAR(50) NOT NULL,
    category VARCHAR(30) DEFAULT 'general',
    question TEXT,
    status VARCHAR(30) NOT NULL,
    is_free_quota BOOLEAN DEFAULT FALSE,
    coins_charged INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS reading_drawn_cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES reading_sessions(id) ON DELETE CASCADE,
    card_id INT NOT NULL,
    position_index INT NOT NULL,
    orientation VARCHAR(10) NOT NULL,
    drawn_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS reading_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID UNIQUE NOT NULL REFERENCES reading_sessions(id) ON DELETE CASCADE,
    overall_summary TEXT NOT NULL,
    position_details JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_reading_sessions_user ON reading_sessions(user_id, created_at DESC);
