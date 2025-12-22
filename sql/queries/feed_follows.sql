-- name: CreateFeedFollows :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT inserted_feed_follow.*, feeds.name AS feed_name, users.name as user_name
FROM inserted_feed_follow
INNER JOIN users ON inserted_feed_follow.user_id = users.id
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
;

-- name: GetFeedFollowsForUser :many
SELECT ff.*, f.name, u.name as user_name FROM feed_follows AS ff
LEFT JOIN users as u on u.id = ff.user_id
LEFT JOIN feeds as f on f.id = ff.feed_id
WHERE ff.user_id = $1;

-- name: UnfollowFeed :exec
DELETE FROM feed_follows
WHERE $1 = feed_follows.user_id
AND $2 = feed_follows.feed_id;