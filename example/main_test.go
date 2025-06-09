package main

import "testing"

func Test_exampleFunc(t *testing.T) {
	t.Run("exampleFunc", func(t *testing.T) {
		exampleFunc(1)
	})
}

func Test_exampleFunc2(t *testing.T) {
	t.Run("exampleFunc2", func(t *testing.T) {
		exampleFunc2(3)
	})
}
