-- +goose Up
-- +goose StatementBegin
CREATE TABLE feature_flags (
    id UUID DEFAULT uuid_generate_v4() NOT NULL
        CONSTRAINT feature_flags_pk PRIMARY KEY,
    owner_id UUID NOT NULL
        CONSTRAINT feature_flags_owner_fk REFERENCES users(id) ON DELETE CASCADE,    
    name VARCHAR(255) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT FALSE,
    description TEXT,
    created_at    TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at    TIMESTAMP DEFAULT NOW() NOT NULL
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feature_flags;
-- +goose StatementEnd