package mapper_test

import (
	"errors"
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
	if user.FirstName == "" {
		return UserDTO{}, errors.New("first name cannot be empty")
	}

	return UserDTO{
		FullName: user.FirstName + " " + user.LastName,
		Age:      user.Age,
		IsAdult:  user.Age >= 18,
	}, nil
}

func ExampleMap() {
	mapper.RegisterMapper[User, UserDTO](MapUserToDTO)

	dto, err := mapper.Map[User, UserDTO](User{FirstName: "John", LastName: "Doe", Age: 25})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", dto)
	// Output: {FullName:John Doe Age:25 IsAdult:true}
}

func ExampleMap_notRegistered() {
	_, err := mapper.Map[UserDTO, User](UserDTO{})

	fmt.Println(errors.Is(err, mapper.ErrMapperNotFound))
	// Output: true
}

func ExampleMapList() {
	mapper.RegisterMapper[User, UserDTO](MapUserToDTO)

	users := []User{
		{FirstName: "Alice", LastName: "Johnson", Age: 28},
		{FirstName: "Bob", LastName: "Brown", Age: 15},
	}

	userDTOS, err := mapper.MapList[User, UserDTO](users)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, dto := range userDTOS {
		fmt.Printf("%+v\n", dto)
	}
	// Output:
	// {FullName:Alice Johnson Age:28 IsAdult:true}
	// {FullName:Bob Brown Age:15 IsAdult:false}
}

func ExampleMapList_error() {
	mapper.RegisterMapper[User, UserDTO](MapUserToDTO)

	_, err := mapper.MapList[User, UserDTO]([]User{{FirstName: "Alice"}, {}})

	fmt.Println(err)
	// Output: mapping element 1: first name cannot be empty
}

func ExampleUnregisterMapper() {
	mapper.RegisterMapper[User, UserDTO](MapUserToDTO)
	fmt.Println(mapper.HasMapping[User, UserDTO]())

	fmt.Println(mapper.UnregisterMapper[User, UserDTO]())
	fmt.Println(mapper.HasMapping[User, UserDTO]())
	// Output:
	// true
	// true
	// false
}
