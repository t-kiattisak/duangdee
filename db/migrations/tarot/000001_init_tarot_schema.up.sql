-- Up Migration for tarot_db
CREATE TABLE IF NOT EXISTS tarot_cards (
    id INT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    arcana_type VARCHAR(20) NOT NULL,
    suit VARCHAR(20),
    number INT NOT NULL,
    element VARCHAR(20),
    image_url TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS card_meanings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id INT NOT NULL REFERENCES tarot_cards(id),
    orientation VARCHAR(10) NOT NULL,
    category VARCHAR(30) NOT NULL,
    meaning TEXT NOT NULL,
    keywords TEXT[] NOT NULL DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS tarot_spreads (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    card_count INT NOT NULL,
    coin_cost INT NOT NULL DEFAULT 0,
    positions_json JSONB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_card_meanings_lookup ON card_meanings(card_id, orientation, category);
