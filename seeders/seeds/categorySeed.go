package seeds

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/a4anthony/go-store/config"
	"github.com/a4anthony/go-store/internal/database"
	"github.com/a4anthony/go-store/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"math/rand"
	"strings"
)

func CreateCategory() error {

	name := gofakeit.ProductCategory()
	description := sql.NullString{}
	description.String = "This is the description of " + name + " category."
	description.Valid = true

	status := rand.Intn(2) == 1
	// check if category with is_active = true exceeds 5
	var count int64
	categories, _ := config.DB.GetCategories(context.Background())
	// count the number of active categories
	for range categories {
		count++
	}
	if count >= 2 {
		status = true
	}
	//capitalizeName := strings.ToUpper(name[:1]) + name[1:]
	capitalizeName := utils.CapitalizeEachWord(name)
	_, err := config.DB.CreateCategory(context.Background(), database.CreateCategoryParams{
		ID:          uuid.New(),
		Name:        capitalizeName,
		Description: description,
		IsActive:    status,
		Slug:        strings.ToLower(strings.ReplaceAll(name, " ", "-")),
	})
	if err != nil {
		fmt.Println(err)
		fmt.Println("err")
	}
	//err = SubCreateCategory()
	//if err != nil {
	//	return err
	//} else {
	//	log.Fatalf("Error creating sub category")
	//}
	return nil
}

func SubCreateCategory() error {
	//var isActive sql.NullBool
	//

	categories, _ := config.DB.GetCategories(context.Background())
	fmt.Println("categories")

	for _, category := range categories {
		// create 5 sub categories for each category
		for j := 0; j < 5; j++ {
			name := "Sub category " + fmt.Sprintf("%d", j+1)
			status := rand.Intn(2) == 1

			_, err := config.DB.CreateSubCategory(context.Background(), database.CreateSubCategoryParams{
				ID:          uuid.New(),
				Name:        name,
				Description: sql.NullString{String: "This is the description of sub category.", Valid: true},
				IsActive:    status,
				Slug:        strings.ToLower(strings.ReplaceAll(name, " ", "-")),
				CategoryID:  category.ID,
			})
			if err != nil {
				fmt.Println(err)
				fmt.Println("err")
			}
		}
	}

	return nil

}
