CREATE TABLE maximum_spends (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    daily_limit DECIMAL(12,2),
    monthly_limit DECIMAL(12,2),
    yearly_limit DECIMAL(12,2),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
