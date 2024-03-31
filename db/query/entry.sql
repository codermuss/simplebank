-- CREATE TABLE "entries" (
--   "id" bigserial PRIMARY KEY,
--   "account_id" bigint NOT NULL,
--   "amount" bigint NOT NULL,
--   "created_at" timestamptz NOT NULL DEFAULT (now())
-- );

-- name: CreateEntry :one
INSERT INTO entries (account_id, amount) VALUES ($1,$2) RETURNING *;