package main

import (
	"fmt"

	"github.com/github-trending/github-trending"
)

func main() {
	repos, _ := trending.Today().Repos()

	for index, repo := range repos {
		fmt.Println(index + 1, repo.Title)
	}

	// fmt.Println(trending.Today().Repos())
	// fmt.Println(trending.Today().WithLanguage("go").Repos())
	// fmt.Println(trending.ThisWeek().Repos())
	// fmt.Println(trending.ThisMonth().Repos())
	// fmt.Println(trending.Repos())
}
