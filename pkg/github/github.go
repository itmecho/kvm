package github

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
)

// Github handles calls to the github api
type Github struct {
	client *http.Client
	apiURL string
}

// New returns a new instance of Github
func New(httpClient *http.Client) Github {
	return Github{
		apiURL: "https://api.github.com",
		client: httpClient,
	}
}

// Release stores the required data from a release
type Release struct {
	Name string `json:"name"`
}

// GetReleases calls the github api and retrieves a list of release names for the given repository
func (gh Github) GetReleases(owner, repo, filter string) ([]Release, error) {
	var releases []Release

	// TODO improve this
	url := strings.Join([]string{gh.apiURL, "repos", owner, repo, "releases"}, "/")
	resp, err := gh.client.Get(url)
	if err != nil {
		return releases, err
	}

	if err = json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return releases, err
	}

	sort.Slice(releases, func(i, j int) bool {
		return releases[i].Name < releases[j].Name
	})

	if filter != "" {
		var filtered []Release

		for _, r := range releases {
			if strings.Contains(r.Name, filter) {
				filtered = append(filtered, r)
			}
		}

		return filtered, nil
	}

	return releases, nil
}
