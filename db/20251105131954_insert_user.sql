-- +goose Up
-- +goose StatementBegin
INSERT INTO users(full_name, password, phone, email) 
VALUES('John Doe','123','+8985688822','john@dmail.com');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE users;
-- +goose StatementEnd
