-- migrate:up
CREATE TABLE Cache (
    key text PRIMARY KEY,
    data bytea NOT NULL,
    expires timestamptz NOT NULL
);

-- migrate:down
DROP TABLE Cache;