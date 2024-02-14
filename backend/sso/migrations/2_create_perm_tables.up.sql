CREATE TABLE IF NOT EXISTS permissions
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS user_permissions
(
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    permission_id INT REFERENCES permissions(id)
);

