-- +goose Up
-- +goose StatementBegin
CREATE TABLE points (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    point BIGINT
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE points;

-- +goose StatementEnd