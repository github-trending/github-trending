package trending

// Repos returns a slice of repositories with default trending client.
func Repos() ([]Repository, error) {
	return New().Repos()
}

// Today sets time span to daily with default trending client.
func Today() *Trending {
	return New().Since("daily")
}

// ThisWeek sets time span to weekly with default trending client.
func ThisWeek() *Trending {
	return New().Since("weekly")
}

// ThisMonth sets time span to monthly with default trending client.
func ThisMonth() *Trending {
	return New().Since("monthly")
}
