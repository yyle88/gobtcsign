package utils

import "fmt"

func MustDone(erx error) {
	if erx != nil {
		panic(erx)
	}
}

func MustEquals[T comparable](want, data T) {
	if want != data {
		fmt.Println("want:", want)
		fmt.Println("data:", data)
		panic("wrong")
	}
}
