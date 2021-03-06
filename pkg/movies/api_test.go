package movies

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

)

func TestSearchMovies(t *testing.T) {
	cases := []struct {
		name                string
		mockResponseBody    string
		expectedMovies      []Movie
		expectedErrorString string
		pagina string
		urlBase string
	}{
		{
			name:             "RegularCase",
			mockResponseBody: `{"Search":[{"Title":"Star Wars: A New Hope","Year":"1977"},{"Title":"Star Wars: The Empire Strikes Back","Year":"1980"}]}`,
			expectedMovies: []Movie{
				{Title: "Star Wars: A New Hope", Year: "1977"},
				{Title: "Star Wars: The Empire Strikes Back", Year: "1980"},
			},
			expectedErrorString: "",
			urlBase:"http://example.com/",
		},
		{
			name:             "RegularPageCase",
			mockResponseBody: `{"Search":[{"Title":"Star Wars: A New Hope","Year":"1977"},{"Title":"Star Wars: The Empire Strikes Back","Year":"1980"}]}`,
			expectedMovies: []Movie{
				{Title: "Star Wars: A New Hope", Year: "1977"},
				{Title: "Star Wars: The Empire Strikes Back", Year: "1980"},
			},
			expectedErrorString: "",
			pagina:"2",
			urlBase:"http://example.com/",
		},
		{
			name:             "ErrorUrlCase",
			mockResponseBody: `{"Search"`,
			expectedMovies: []Movie(nil),
			expectedErrorString: "Get \"http://example2.com/?apikey=mock-api-key&s=star+wars&type=movie\": no responder found",
			urlBase:"http://example2.com/",
		},
	}

	

	for _, c := range cases {

		searcher := &APIMovieSearcher{
			URL:    c.urlBase,
			APIKey: "mock-api-key",
		}

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
			if(c.pagina!=""){
				searchQuery["p"] = c.pagina
			}

			//actualMovies, actualError := searcher.SearchMovies("star wars")
			actualMovies, actualError := searcher.SearchMovies(searchQuery)
			assert.EqualValues(t, c.expectedMovies, actualMovies)


			if c.expectedErrorString == "" {
				assert.NoError(t, actualError)
			} else {
				assert.EqualError(t, actualError, c.expectedErrorString)
			}
		})

		
	}


}
