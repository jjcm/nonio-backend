-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD `current_period_end` TIMESTAMP,ADD `next_payout` TIMESTAMP,ADD `last_payout` TIMESTAMP;
-- +goose StatementEnd
