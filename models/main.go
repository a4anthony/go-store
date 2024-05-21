package models

import (
	"database/sql"
	"github.com/a4anthony/go-store/internal/database"
	"github.com/a4anthony/go-store/utils"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID              uuid.UUID  `json:"id"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	Phone           string     `json:"phone"`
	Email           string     `json:"email"`
	Password        string     `json:"-"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	Avatar          string     `json:"avatar"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type Category struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SubCategory struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description sql.NullString `json:"description"`
	IsActive    bool           `json:"is_active"`
	CategoryID  uuid.UUID      `json:"category_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func DbUserToUser(dbUser database.User) User {
	background := utils.GenerateColorByFirstLetter(dbUser.FirstName)
	avatar := "https://ui-avatars.com/api/?name=" + dbUser.FirstName + "+" + dbUser.LastName + "&size=256&background=" + background + "&color=fff"
	return User{
		ID:              dbUser.ID,
		FirstName:       dbUser.FirstName,
		LastName:        dbUser.LastName,
		Phone:           dbUser.Phone,
		Email:           dbUser.Email,
		Password:        dbUser.Password,
		EmailVerifiedAt: nullTimeToTimePtr(dbUser.EmailVerifiedAt),
		CreatedAt:       dbUser.CreatedAt,
		UpdatedAt:       dbUser.UpdatedAt,
		Avatar:          avatar,
	}

}

func DbCategoryToCategory(category database.Category) Category {
	return Category{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: nullStringToStringPtr(category.Description),
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

func DbSubCategoryToSubCategory(category database.SubCategory) SubCategory {
	return SubCategory{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		IsActive:    category.IsActive,
		CategoryID:  category.CategoryID,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

func DbCategoriesListToCategoriesList(categories []database.Category) []Category {
	var categoriesList []Category
	for _, category := range categories {
		categoriesList = append(categoriesList, DbCategoryToCategory(category))
	}
	return categoriesList
}

func nullTimeToTimePtr(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func nullStringToStringPtr(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

type Sorting interface {
	SortByCreatedAt()
}

func (u *User) SortByCreatedAt() {

}
