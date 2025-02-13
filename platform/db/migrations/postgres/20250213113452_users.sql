-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id            UUID DEFAULT uuid_generate_v4() NOT NULL
        CONSTRAINT users_pk PRIMARY KEY,
    name          VARCHAR NOT NULL,
    email         VARCHAR NOT NULL
        CONSTRAINT users_email_unique UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at    TIMESTAMP DEFAULT NOW() NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
