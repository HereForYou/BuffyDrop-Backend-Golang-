package utils

import "fmt"

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