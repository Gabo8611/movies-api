package movies

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"sort"
)

func TestMovies(t *testing.T) {
	cases := []struct {
		name                string
		mockResponseBody    string
		expectedMovies      []Movie
		expectedErrorString string
	}{
		{
			name:             "SortCase",
			mockResponseBody: `{"Search":[{"Title":"Star Wars: The Empire Strikes Back","Year":"1980"},{"Title":"Star Wars: A New Hope","Year":"1977"}]}`,
			expectedMovies: []Movie{
				{Title: "Star Wars: A New Hope", Year: "1977"},
				{Title: "Star Wars: The Empire Strikes Back", Year: "1980"},
			},
			expectedErrorString: "",
		},
		{
			name:             "SorYearCase",
			mockResponseBody: `{"Search":[{"Title":"Star Wars: The Empire Strikes Back","Year":"1980"},{"Title":"Star Wars: The Empire Strikes Back","Year":"1977"}]}`,
			expectedMovies: []Movie{
				{Title: "Star Wars: The Empire Strikes Back", Year: "1977"},
				{Title: "Star Wars: The Empire Strikes Back", Year: "1980"},
			},
			expectedErrorString: "",
		},
	}

	searcher := &APIMovieSearcher{
		URL:    "http://example.com/",
		APIKey: "mock-api-key",
	}

	for _, c := range cases {
		// register http mock
		httpmock.RegisterResponder(
			"GET",
			"http://example.com/",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(200, c.mockResponseBody), nil
			},
		)
		
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		// run test
		t.Run(c.name, func(t *testing.T) {

			searchQuery := make(map[string]interface{})
			searchQuery["q"] = "star wars"

			//actualMovies, actualError := searcher.SearchMovies("star wars")
			actualMovies, actualError := searcher.SearchMovies(searchQuery)

			movies_sort := Movies(actualMovies)
			sort.Sort(Movies(movies_sort))

			assert.EqualValues(t, c.expectedMovies, actualMovies)

			if c.expectedErrorString == "" {
				assert.NoError(t, actualError)
			} else {
				assert.EqualError(t, actualError, c.expectedErrorString)
			}
		})
	}


}
