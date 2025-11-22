-- Cache queries using hstore for key-value data storage
-- name: GetCache :one
SELECT data, expires FROM cache
WHERE key = $1 LIMIT 1;

-- name: InsertCache :exec
INSERT INTO cache (key, data, expires)
VALUES ($1, $2, $3)
ON CONFLICT(key) DO UPDATE SET
    data = excluded.data,
    expires = excluded.expires;

-- name: DeleteOld :exec
DELETE FROM cache WHERE expires <= $1;