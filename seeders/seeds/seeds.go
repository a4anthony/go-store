package seeds

import (
	"fmt"
	"github.com/a4anthony/go-store/seeders"
	"github.com/brianvoe/gofakeit/v7"
)

func All() []seeders.Seed {
	var output []seeders.Seed
	output = append(output, userSeed()...)
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
