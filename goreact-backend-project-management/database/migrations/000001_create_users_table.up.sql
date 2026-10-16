CREATE TABLE users
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(255)        NOT NULL,
    email      VARCHAR(255) UNIQUE NOT NULL,
    password   TEXT                NOT NULL,
    role       VARCHAR(50)         NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ                  DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ                  DEFAULT CURRENT_TIMESTAMP
);
