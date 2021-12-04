-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD `express_account_id` VARCHAR(64) NOT NULL DEFAULT "" AFTER stripe_customer_id;
-- +goose StatementEnd
