-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE `posts_tags` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `post_id` int unsigned NOT NULL DEFAULT 0,
  `tag_id` int unsigned NOT NULL DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT NOW(),
  PRIMARY KEY (`id`),
  CONSTRAINT `post_tags` UNIQUE (post_id, tag_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
