CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    role TEXT DEFAULT 'user',
    status TEXT DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    total NUMERIC,
    status TEXT DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    price NUMERIC NOT NULL,
    stock INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id),
    product_id INT REFERENCES products(id),
    quantity INT NOT NULL,
    unit_price NUMERIC NOT NULL
);

CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    action TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id INT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Insert sample data for realistic EXPLAIN ANALYZE results
INSERT INTO users (name, email, role, status) VALUES
    ('Alice', 'alice@test.com', 'admin', 'active'),
    ('Bob', 'bob@test.com', 'user', 'active'),
    ('Charlie', 'charlie@test.com', 'user', 'inactive'),
    ('Diana', 'diana@test.com', 'manager', 'active');

INSERT INTO products (name, category, price, stock) VALUES
    ('Laptop', 'electronics', 999.99, 50),
    ('Phone', 'electronics', 599.99, 200),
    ('Desk', 'furniture', 249.99, 30),
    ('Chair', 'furniture', 149.99, 100);

INSERT INTO orders (user_id, total, status) VALUES
    (1, 1599.98, 'completed'),
    (2, 249.99, 'pending'),
    (1, 149.99, 'completed'),
    (3, 599.99, 'cancelled');

INSERT INTO order_items (order_id, product_id, quantity, unit_price) VALUES
    (1, 1, 1, 999.99),
    (1, 2, 1, 599.99),
    (2, 3, 1, 249.99),
    (3, 4, 1, 149.99),
    (4, 2, 1, 599.99);
