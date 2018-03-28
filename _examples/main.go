package main

import (
	"fmt"

	"github.com/github-trending/github-trending"
)

func main() {
	fmt.Println(trending.Today().Repos())
	// fmt.Println(trending.ThisWeek().Repos())
	// fmt.Println(trending.ThisMonth().Repos())
	// fmt.Println(trending.Repos())
}
