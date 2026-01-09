-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeedsUsers :many
SELECT feeds.id, feeds.created_at, feeds.updated_at, feeds.name as feed_name, feeds.url, feeds.user_id, users.name as user_name
FROM feeds
LEFT JOIN users
ON feeds.user_id = users.id
ORDER BY feeds.updated_at DESC;

-- name: GetFeedURL :one
SELECT id, created_at, updated_at, name, url, user_id
FROM feeds
WHERE url = $1
LIMIT 1;