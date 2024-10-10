package utils

import (
	"fmt"
	"go-test/models"
)

func SayHello(name string) {
	fmt.Println("Hello, ", name)
}

func FindEvens(nums int)  {
	for i := 0; i < nums ; i++ {
		if i%2 ==0 {
			fmt.Println(i)
		}
	}
}

func HasFriendWithId (friends []models.Friend, id string) bool {
	for _, friend := range friends {
		if friend.Id == id {
			return true
		}
	}
	return false
}