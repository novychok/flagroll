-- +goose Up
-- +goose StatementBegin
create table api_keys (
    id uuid default uuid_generate_v4() not null
            constraint api_keys_pk
                primary key,
    owner_id uuid not null
        constraint api_keys_fk
         references users(id) on delete cascade,
    token_id uuid not null,
    token_hash varchar not null,
    created_at timestamp not null default now(),
    expires_at timestamp not null
);

CREATE UNIQUE INDEX idx_api_keys_token_id ON api_keys (token_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index idx_api_keys_token_id;

drop table api_keys;
-- +goose StatementEnd
