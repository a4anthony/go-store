package handlers

import (
	"database/sql"
	"fmt"
	"github.com/a4anthony/go-store/config"
	"github.com/a4anthony/go-store/internal/database"
	"github.com/a4anthony/go-store/models"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type GetCategoriesParams struct {
	IsActive bool `json:"is_active"`
}

// Filter represents the filtering criteria
type Filter struct {
	Name     *string
	IsActive *bool
}

type QueriesFilter struct {
	Db database.DBTX
}

func NewFilter(db database.DBTX) *QueriesFilter {
	return &QueriesFilter{Db: db}
}

// ListUsersWithFilters dynamically constructs the SQL query based on the provided filters
func ListUsersWithFilters(ctx *fiber.Ctx, q *QueriesFilter, filters Filter) ([]database.Category, error) {
	baseQuery := "SELECT id, name, is_active FROM categories where 1 = 1"
	var args []interface{}
	var conditions []string

	if filters.Name != nil && *filters.Name != "" {
		conditions = append(conditions, "name = $1")
		args = append(args, *filters.Name)
	}
	if *filters.Name == "" {
		args = append(args, *filters.Name)
	}

	if filters.IsActive != nil {
		conditions = append(conditions, "is_active = $2")
		args = append(args, *filters.IsActive)
	}
	fmt.Println(conditions)
	fmt.Println(args)
	if len(conditions) > 0 {
		baseQuery = fmt.Sprintf("%s AND %s", baseQuery, strings.Join(conditions, " AND "))
	}
	//fmt.Println(baseQuery)
	//baseQuery = "SELECT id, name, is_active FROM categories where 1 = 1 AND is_active = $2"
	//baseQuery = "SELECT * FROM categories WHERE (COALESCE($1::text, '') = '' OR name ILIKE '%' || $1 || '%')  AND (COALESCE($2::boolean, true) = true OR is_active = $2)"
	baseQuery = "SELECT * FROM categories WHERE (COALESCE($1::text, '') = '' OR name ILIKE '%' || $1 || '%')  AND (COALESCE($2::boolean, false) IS NULL OR is_active = $2)"
	//baseQuery = "SELECT id, name, is_active FROM categories where 1 = 1 AND name = $1 AND is_active = $2"
	fmt.Println(args...)
	rows, err := q.Db.QueryContext(ctx.Context(), baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []database.Category
	for rows.Next() {
		var i database.Category
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Slug,
			&i.Description,
			&i.IsActive,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		categories = append(categories, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func GetCategories(c *fiber.Ctx) error {
	//isActiveFilter := c.Query("is_active") == "1" || c.Query("is_active") == "true"
	name := c.Query("name")

	isActive := sql.NullBool{Valid: false}
	if c.Query("is_active") == "1" || c.Query("is_active") == "true" {
		isActive = sql.NullBool{Bool: true, Valid: true}
	} else if c.Query("is_active") == "0" || c.Query("is_active") == "false" {
		isActive = sql.NullBool{Bool: false, Valid: true}
	} else {
		fmt.Println("is_active is not valid")
		isActive = sql.NullBool{Valid: false}
	}

	fmt.Println(isActive)

	//filters := Filter{
	//	Name:     &name,
	//	IsActive: &isActive,
	//}
	//
	//connStr := os.Getenv("DB_URL")
	//dbConn, err := sql.Open("postgres", connStr)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//queries := NewFilter(dbConn)

	//categories, categoriesErr := ListUsersWithFilters(c, queries, filters)
	categories, categoriesErr := config.DB.GetCategories(c.Context(), database.GetCategoriesParams{
		Column1: name,
		Column2: isActive,
		//Column2: sql.NullBool{Valid: false},
		//Column2: sql.NullBool{Bool: isActive, Valid: true},
	})

	if categoriesErr != nil {
		fmt.Println(categoriesErr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting categories",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		//"data": categories,
		"data":    models.DbCategoriesListToCategoriesList(categories),
		"message": "Categories retrieved successfully.",
	})
}
