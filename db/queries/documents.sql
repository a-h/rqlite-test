-- name: DocumentsInsert :exec
INSERT INTO "documents" ("name", "content") VALUES (?, ?);

-- name: DocumentsSelectMany :many
SELECT "id", "name", "content" FROM "documents";

-- name: DocumentsSelectOneByID :one
SELECT "id", "name", "content" FROM "documents" WHERE "id" = ?;
