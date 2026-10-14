package main

import (
	"fmt"

	"github.com/mdnishat27/go-mapper"
)

type User struct {
	FirstName string
	LastName  string
	Age       int
}

type UserDTO struct {
	FullName string
	Age      int
	IsAdult  bool
}

func MapUserToDTO(user User) (UserDTO, error) {
	return UserDTO{
		FullName: user.FirstName + " " + user.LastName,
		Age:      user.Age,
		IsAdult:  user.Age >= 18,
	}, nil
}

func main() {
	// Register the mapping function
	mapper.RegisterMapper[User, UserDTO](MapUserToDTO)

	// Create test users
	users := []User{
		{FirstName: "Alex", LastName: "Russo", Age: 25},
		{FirstName: "Ben", LastName: "Tennyson", Age: 10},
		{FirstName: "Chloe", LastName: "Thunderman", Age: 15},
	}

	fmt.Println("=== Mapping Users to DTOs ===")

	// Map each user
	for i, user := range users {
		dto, err := mapper.Map[User, UserDTO](user)
		if err != nil {
			fmt.Printf("Error mapping user %d: %v\n", i+1, err)
			continue
		}

		fmt.Printf("User %d: %s -> %s (Age: %d, Adult: %t)\n",
			i+1,
			user.FirstName+" "+user.LastName,
			dto.FullName,
			dto.Age,
			dto.IsAdult,
		)
	}

	// Demonstrate list mapping
	fmt.Println("\n=== Bulk Mapping ===")
	userDTOS, err := mapper.MapList[User, UserDTO](users)
	if err != nil {
		fmt.Printf("Error in bulk mapping: %v\n", err)
		return
	}

	fmt.Printf("Mapped %d users to %d DTOs\n", len(users), len(userDTOS))

	// Demonstrate mapping management
	fmt.Println("\n=== Mapping Management ===")

	// Check if mapping exists
	if mapper.HasMapping[User, UserDTO]() {
		fmt.Println("✓ User->UserDTO mapping exists")
	}

	// Remove mapping
	removed := mapper.UnregisterMapper[User, UserDTO]()
	if removed {
		fmt.Println("✓ User->UserDTO mapping removed")
	}

	// Check again
	if !mapper.HasMapping[User, UserDTO]() {
		fmt.Println("✓ User->UserDTO mapping no longer exists")
	}

	fmt.Println("\n=== Example Complete ===")
}
