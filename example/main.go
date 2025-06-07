package main

import "math/rand"

func main() {
	value := rand.Intn(100)
	exampleFunc(value)
	exampleFunc2(value)
}

func exampleFunc(value int) {
	if value%2 == 0 {
		println("Even")
	} else {
		println("Odd")
	}
}

func exampleFunc2(value int) {
	if value%3 == 0 {
		println("Divisible by 3")
	} else {
		println("Not divisible by 3")
	}
}
