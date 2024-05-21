package seeds

import (
	"fmt"
	"github.com/a4anthony/go-store/seeders"
	"github.com/brianvoe/gofakeit/v7"
)

func All() []seeders.Seed {
	var output []seeders.Seed
	output = append(output, userSeed()...)
	output = append(output, categorySeed()...)
	output = append(output, subCategorySeed()...)
	return output
}

func userSeed() []seeders.Seed {
	fmt.Println("Creating users...")
	// an array of string Jane, John, Susan
	var names []string
	for i := 0; i < 11; i++ {
		names = append(names, gofakeit.Person().FirstName)
	}
	seeds := make([]seeders.Seed, len(names))

	for i, name := range names {
		email := gofakeit.Email()
		if i == 0 {
			email = "johndoe@email.com"
		}
		seeds[i] = seeders.Seed{
			Name: "Create" + name,
			Run: func() error {
				err := CreateUser(name, gofakeit.Person().LastName, "+"+gofakeit.Phone(), email)
				return err
			},
		}
	}

	return seeds
}

func categorySeed() []seeders.Seed {
	fmt.Println("Creating categories...")
	seeds := make([]seeders.Seed, 10)
	// create 10 categories
	for i := 0; i < 10; i++ {
		seeds[i] = seeders.Seed{
			Name: "CreateCategory",
			Run: func() error {
				err := CreateCategory()
				return err
			},
		}
	}
	return seeds
}

func subCategorySeed() []seeders.Seed {
	fmt.Println("Creating sub categories...")
	seeds := make([]seeders.Seed, 1)
	seeds[0] = seeders.Seed{
		Name: "CreateSubCategory",
		Run: func() error {
			err := SubCreateCategory()
			return err
		},
	}
	return seeds
}
