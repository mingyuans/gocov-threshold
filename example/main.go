package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	value := rand.Intn(100)
	exampleFunc(value)
	exampleFunc2(value)
	exampleFunc3(value)
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
		// gocover:ignore
		println("Not divisible by 3")
	}
}

func exampleFunc3(value int) {
	// init a mutex and lock
	var mu = sync.Mutex{}
	mu.Lock()
	fmt.Println("Mutex locked")
	defer mu.Unlock()
	if value%5 == 0 {
		fmt.Println("Divisible by 5")
	}
}

func exampleFunc4(value int) {
	// init a mutex and lock
	var mu = sync.Mutex{}
	mu.Lock()
	fmt.Println("Mutex locked")
	defer mu.Unlock()
	fmt.Println("Mutex unlocked")
}
