// Package trending tracks the most popular GitHub repos.
package trending

import (
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	defaultBaseURL   = "https://github.com"
	trendingPathname = "/trending"
)

// Repository reflects a single trending repository.
type Repository struct {
	Title           string
	Owner           string
	Name            string
	Description     string
	Language        string
	Stars           int
	AdditionalStars int
	URL             *url.URL
}

type Trending struct {
	timeSpan string
	BaseURL  *url.URL
	Client   *http.Client
}

// New returns a new trending point.
func New() *Trending {
	return NewWithClient(http.DefaultClient)
}

// NewWithClient allows providing a custom http.Client to use for fetching data.
func NewWithClient(c *http.Client) *Trending {
	baseURL, _ := url.Parse(defaultBaseURL)

	t := Trending{
		BaseURL: baseURL,
		Client:  c,
	}

	return &t
}

// Since allows adjusting time span.
// string `ts` should be one of daily, weekly or monthly.
func (t *Trending) Since(ts string) *Trending {
	t.timeSpan = ts

	return t
}

// Repos returns a slice of fetched trending repositories.
func (t *Trending) Repos() ([]Repository, error) {
	var repositories []Repository
	var ts string

	if t.timeSpan != "" {
		ts = t.timeSpan
	} else {
		ts = "daily"
	}

	u, err := t.formatURL(ts, "")

	if err != nil {
		return repositories, err
	}

	doc, err := t.request(u)

	if err != nil {
		return repositories, err
	}

	repositories = t.separateRepos(doc)

	return repositories, nil
}

func (t *Trending) formatURL(since string, language string) (string, error) {
	pathname := trendingPathname

	if language != "" {
		pathname += "/" + language
	}

	path, err := url.Parse(pathname)

	if err != nil {
		return "", err
	}

	u := t.BaseURL.ResolveReference(path)

	query := u.Query()

	query.Add("since", since)

	u.RawQuery = query.Encode()

	return u.String(), nil
}

func (t *Trending) request(u string) (*goquery.Document, error) {
	response, err := t.Client.Get(u)
	defer response.Body.Close()

	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		return nil, err
	}

	return doc, nil
}

// separateRepos returns a slice of repositories that separated from html.
func (t *Trending) separateRepos(doc *goquery.Document) []Repository {
	var repositories []Repository

	doc.Find("ol.repo-list li").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h3 a").Text())
		splittedName := strings.Split(title, "/")
		owner := strings.TrimSpace(splittedName[0])
		repositoryName := strings.TrimSpace(splittedName[1])

		description := strings.TrimSpace(s.Find(".py-1 p").Text())
		language := strings.TrimSpace(s.Find("div.f6 > span").Eq(0).Text())

		starsStr := strings.TrimSpace(s.Find("div.f6 a").First().Text())
		stars, _ := strconv.Atoi(strings.Replace(starsStr, ",", "", -1))

		additionalStarsText := strings.TrimSpace(s.Find("div.f6 span.float-sm-right").Text())
		additionalStarsStr := regexp.MustCompile("[0-9]+").FindString(additionalStarsText)
		additionalStars, _ := strconv.Atoi(additionalStarsStr)

		repositoryURLStr, exists := s.Find("h3 a").First().Attr("href")

		var repositoryURL *url.URL

		if exists {
			repositoryURL, _ = t.BaseURL.Parse(repositoryURLStr)
		} else {
			repositoryURL, _ = t.BaseURL.Parse(strings.Replace(title, " ", "", -1))
		}

		repo := Repository{
			Title:           title,
			Owner:           owner,
			Name:            repositoryName,
			Description:     description,
			Language:        language,
			Stars:           stars,
			AdditionalStars: additionalStars,
			URL:             repositoryURL,
		}

		repositories = append(repositories, repo)
	})

	return repositories
}
