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

-- name: GetFeeds :many
select * FROM feeds;

-- name: GetFeedByUrl :one
select * from feeds
where feeds.url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = $2, last_fetched_at = $2
WHERE feeds.id = $1;

-- name: GetNextFeedToFetch :one
SELECT f.* FROM feeds as f
INNER JOIN feed_follows as ff on ff.feed_id = f.id
WHERE ff.user_id = $1
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
