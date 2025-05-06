-- +goose Up
CREATE TABLE metrics
(
    name  TEXT PRIMARY KEY,
    value DOUBLE PRECISION NOT NULL
);

-- +goose Down
DROP TABLE metrics
