package utils

import "fmt"

func MustDone(err error) {
	if err != nil {
		panic(err)
	}
}

func MustEquals[T comparable](want, data T) {
	if want != data {
		fmt.Println("want:", want)
		fmt.Println("data:", data)
		panic("wrong")
	}
}
