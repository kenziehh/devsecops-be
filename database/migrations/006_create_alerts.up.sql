CREATE TABLE alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    message VARCHAR(255),
    type VARCHAR(20), -- daily, monthly, yearly
    triggered_at TIMESTAMP DEFAULT NOW()
);
