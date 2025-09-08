CREATE TYPE transaction_type AS ENUM ('income', 'expense');

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    category_id UUID REFERENCES categories(id),
    type transaction_type NOT NULL,
    period VARCHAR(20),
    amount DECIMAL(12,2),
    note TEXT,
    date DATE,
    proof_file VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
