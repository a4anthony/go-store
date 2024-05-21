-- name: CreateCategory :one
INSERT INTO categories (id, created_at, updated_at, name, slug, description, is_active)
VALUES ($1, NOW(), NOW(), $2, $3, $4, $5) RETURNING *;


-- -- name: GetCategories :many
-- -- sql query to get all categories and then it will sort the result by name
-- SELECT * FROM categories ORDER BY name ASC;


-- name: GetCategorySubCategories :many
-- sql query to get all subcategories of a category
SELECT * FROM sub_categories WHERE category_id = $1;


-- name: GetCategoriesByStatus :many
-- sql query to get all categories by status
SELECT * FROM categories WHERE is_active = $1 ORDER BY name ASC;



-- name: GetCategories :many
-- SELECT *
-- FROM categories
-- WHERE (COALESCE($1::text, '') = '' OR name ILIKE '%' || $1 || '%')
--   AND (is_active = COALESCE($2::bool, is_active));

SELECT *
FROM categories
WHERE (COALESCE($1::text, '') = '' OR name ILIKE '%' || $1 || '%')
   AND ($2::bool IS NULL OR is_active = $2);


-- SELECT *
-- FROM categories
-- WHERE (COALESCE($1::text, '') = '' OR name ILIKE '%' || $1 || '%')
--   AND ($2 IS NULL OR is_active = $2);