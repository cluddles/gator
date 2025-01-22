-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
JOIN feeds on feeds.id = inserted_feed_follow.feed_id
JOIN users on users.id = inserted_feed_follow.user_id;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE user_id = $1 and feed_id = $2;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, feeds.name AS feed_name, users.name as user_name
FROM feed_follows
JOIN feeds on feeds.id = feed_follows.feed_id
JOIN users on users.id = feed_follows.user_id
WHERE feed_follows.user_id = $1;
