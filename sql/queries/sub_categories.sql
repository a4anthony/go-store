-- name: CreateSubCategory :one
INSERT INTO sub_categories (id, created_at, updated_at, name, slug, description, is_active, category_id)
VALUES ($1, NOW(), NOW(), $2, $3, $4, $5, $6) RETURNING *;