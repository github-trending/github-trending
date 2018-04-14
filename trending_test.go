package trending

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Trending
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client = New()

	u, _ := url.Parse(server.URL)

	client.BaseURL = u
}

func shutdown() {
	server.Close()
}

func newTrending(timeSpan, language string) *Trending {
	trending := New()

	u, _ := url.Parse(defaultBaseURL)

	trending.BaseURL = u
	trending.Client = http.DefaultClient
	trending.timeSpan = timeSpan
	trending.language = language

	return trending
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Trending
	}{
		{"default", newTrending("", "")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWithClient(t *testing.T) {
	type args struct {
		c *http.Client
	}

	tests := []struct {
		name string
		args args
		want *Trending
	}{
		{"default", args{http.DefaultClient}, newTrending("", "")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWithClient(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWithClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrending_Since(t *testing.T) {
	BaseURL, _ := url.Parse(defaultBaseURL)

	type fields struct {
		timeSpan string
		language string
		BaseURL  *url.URL
		Client   *http.Client
	}

	type args struct {
		ts string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Trending
	}{
		{"today", fields{"", "", BaseURL, http.DefaultClient}, args{"daily"}, newTrending("daily", "")},
		{"this week", fields{"", "", BaseURL, http.DefaultClient}, args{"weekly"}, newTrending("weekly", "")},
		{"this month", fields{"", "", BaseURL, http.DefaultClient}, args{"monthly"}, newTrending("monthly", "")},
		{"set to today", fields{"weekly", "", BaseURL, http.DefaultClient}, args{"daily"}, newTrending("daily", "")},
		{"set to this week", fields{"monthly", "", BaseURL, http.DefaultClient}, args{"weekly"}, newTrending("weekly", "")},
		{"set to this month", fields{"daily", "", BaseURL, http.DefaultClient}, args{"monthly"}, newTrending("monthly", "")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trending := &Trending{
				timeSpan: tt.fields.timeSpan,
				language: tt.fields.language,
				BaseURL:  tt.fields.BaseURL,
				Client:   tt.fields.Client,
			}

			if got := trending.Since(tt.args.ts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Trending.Since() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrending_WithLanguage(t *testing.T) {
	BaseURL, _ := url.Parse(defaultBaseURL)

	type fields struct {
		timeSpan string
		language string
		BaseURL  *url.URL
		Client   *http.Client
	}

	type args struct {
		lang string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Trending
	}{
		{"today", fields{"daily", "", BaseURL, http.DefaultClient}, args{"go"}, newTrending("daily", "go")},
		{"this week", fields{"weekly", "", BaseURL, http.DefaultClient}, args{"go"}, newTrending("weekly", "go")},
		{"this month", fields{"monthly", "", BaseURL, http.DefaultClient}, args{"go"}, newTrending("monthly", "go")},
		{"set to golang", fields{"daily", "javascript", BaseURL, http.DefaultClient}, args{"go"}, newTrending("daily", "go")},
		{"set to javascript", fields{"daily", "go", BaseURL, http.DefaultClient}, args{"javascript"}, newTrending("daily", "javascript")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trending := &Trending{
				timeSpan: tt.fields.timeSpan,
				language: tt.fields.language,
				BaseURL:  tt.fields.BaseURL,
				Client:   tt.fields.Client,
			}
			if got := trending.WithLanguage(tt.args.lang); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Trending.WithLanguage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrending_Repos(t *testing.T) {
	setup()
	defer shutdown()

	mux.HandleFunc("/trending", func(w http.ResponseWriter, r *http.Request) {
		contents, err := ioutil.ReadFile("./_local_test_source/trending.html")

		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		fmt.Fprint(w, string(contents))
	})

	repos, err := client.Repos()

	if err != nil {
		t.Errorf("Repos() returned an error: %v", err)
	}

	if len(repos) != 25 {
		t.Errorf("The length of Trending.Repos() = %v, want %v", len(repos), 25)
	}

	got := repos[0]

	u, _ := client.BaseURL.Parse("/schollz/find3")

	want := Repository{
		Title:           "schollz / find3",
		Owner:           "schollz",
		Name:            "find3",
		Description:     "High-precision indoor positioning framework, version 3.",
		Language:        "Go",
		Stars:           814,
		AdditionalStars: 754,
		URL:             u,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Trending.Repos() = %v, want %v", got, want)
	}
}

func TestTrending_FormatURL(t *testing.T) {
	BaseURL, _ := url.Parse(defaultBaseURL)

	type fields struct {
		timeSpan string
		language string
		BaseURL  *url.URL
		Client   *http.Client
	}

	type args struct {
		since    string
		language string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{"today's repos", fields{"daily", "", BaseURL, http.DefaultClient}, args{"daily", ""}, "https://github.com/trending?since=daily", false},
		{"today's repos of go", fields{"daily", "", BaseURL, http.DefaultClient}, args{"daily", "go"}, "https://github.com/trending/go?since=daily", false},
		{"this week's repos of go", fields{"weekly", "javascript", BaseURL, http.DefaultClient}, args{"weekly", "go"}, "https://github.com/trending/go?since=weekly", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trending := &Trending{
				timeSpan: tt.fields.timeSpan,
				language: tt.fields.language,
				BaseURL:  tt.fields.BaseURL,
				Client:   tt.fields.Client,
			}

			got, err := trending.FormatURL(tt.args.since, tt.args.language)

			if (err != nil) != tt.wantErr {
				t.Errorf("Trending.FormatURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Trending.FormatURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
