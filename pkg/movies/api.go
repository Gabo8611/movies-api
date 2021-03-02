package movies

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// APIMovieSearcher is a MovieSearcher implementation using omdbapi
type APIMovieSearcher struct {
	APIKey string
	URL    string
}

type omdbapiResponse struct {
	Search []Movie `json:Search`
}

// SearchMovies searches for a movie
func (s *APIMovieSearcher) SearchMovies(query map[string]interface{}) ([]Movie, error) {

	// call omdbapi
	params := url.Values{}
	params.Add("s", query["q"].(string))
	params.Add("apikey", s.APIKey)
	params.Add("type", "movie")
	if(query["p"]!=nil){
		params.Add("page", query["p"].(string))
	}
	if(query["plot"]!=nil){
		params.Add("plot", query["plot"].(string))
	}


	//params.Add("page", )
	resp, err := http.Get(s.URL + "?" + params.Encode())
	if err != nil {
		return nil, err
	}

	// unmarshall response and get the movie array
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var respStruct omdbapiResponse
	json.Unmarshal(respBody, &respStruct)

	// return result
	return respStruct.Search, nil
}
