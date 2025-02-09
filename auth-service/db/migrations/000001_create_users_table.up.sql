CREATE TABLE IF NOT EXISTS users
(
    id              SERIAL PRIMARY KEY,
    email           VARCHAR(200) NOT NULL UNIQUE,
    password        VARCHAR(200) NOT NULL,
    first_name      VARCHAR(100) NOT NULL,
    last_name       VARCHAR(100) NOT NULL,
    profile_picture VARCHAR      NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);