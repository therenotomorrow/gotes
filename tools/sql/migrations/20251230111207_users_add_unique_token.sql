-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD CONSTRAINT users_token_unique UNIQUE (token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP CONSTRAINT users_token_unique;
-- +goose StatementEnd
