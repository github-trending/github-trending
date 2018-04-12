package trending

// Repos returns a slice of repositories with default trending client.
func Repos() ([]Repository, error) {
	return New().Repos()
}

// Today sets time span to daily with default trending client.
func Today() *Trending {
	return New().Since(today)
}

// ThisWeek sets time span to weekly with default trending client.
func ThisWeek() *Trending {
	return New().Since(thisWeek)
}

// ThisMonth sets time span to monthly with default trending client.
func ThisMonth() *Trending {
	return New().Since(thisMonth)
}
