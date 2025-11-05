-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks (
    id BIGSERIAL PRIMARY KEY,
    task_name VARCHAR(50),
    creator_user_id BIGINT,
    point BIGINT,
    compiler_user_id BIGINT,
    complete_date TIMESTAMP,
    is_complete BOOLEAN
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE tasks;

-- +goose StatementEnd