CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE coins (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    balance INT NOT NULL CHECK (balance >= 0),
    CONSTRAINT positive_balance CHECK (balance >= 0)
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    sender_id UUID REFERENCES users(id) ON DELETE SET NULL,
    receiver_id UUID REFERENCES users(id) ON DELETE SET NULL,
    amount INT NOT NULL CHECK (amount > 0),
    transaction_time TIMESTAMP DEFAULT NOW()
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    price INT NOT NULL CHECK (price >= 0)
);

CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    item_id INT REFERENCES items(id) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK (quantity > 0),
    UNIQUE(user_id, item_id)
);
-- CREATE INDEX idx_inventory_user ON inventory(user_id);


INSERT INTO items (name,price) 
VALUES 
    ('t-shirt','80'),
    ('cup','20'),
    ('book','50'),
    ('pen','10'),
    ('powerbank','200'),
    ('hoody','300'),
    ('umbrella','200'),
    ('socks','10'),
    ('wallet','50'),
    ('pink-hoody','500')
