package main

import "github.com/auvitly/go-tools/nuclear/function"

func Sum(a, b int) int {
	return a + b
}

func main() {
	var (
		old  func(a, b int) int
		impl = func(a, b int) int {
			return 1
		}
	)

	function.Replace(Sum, impl, &old)

	print(Sum(1, 1))
}
