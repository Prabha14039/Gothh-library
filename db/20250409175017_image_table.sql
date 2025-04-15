-- +goose Up
-- +goose StatementBegin

CREATE TABLE images (
    id SERIAL PRIMARY KEY,
    name VARCHAR(1024) NOT NULL,
    url VARCHAR(1024) NOT NULL
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS images;
-- +goose StatementEnd
