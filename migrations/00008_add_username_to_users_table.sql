-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE users ADD `username` VARCHAR(255) NOT NULL DEFAULT "" AFTER email;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
