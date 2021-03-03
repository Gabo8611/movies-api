package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.mpi-internal.com/guillermo-dlsg/movies-api/pkg/movies"

	"strconv"

	// "fmt"
	"sort"
)

func createSearchMoviesHandler(s movies.MovieSearcher) func(http.ResponseWriter, *http.Request) {


	return func(w http.ResponseWriter, req *http.Request) {

		// get parameter from request
		searchQuery := generateQuery(req)
		response := generateResponse(s, searchQuery, w)
		json.NewEncoder(w).Encode(response)
	}
}

func createSearchMoviesSortHandler(s movies.MovieSearcher) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {

		
		searchQuery := generateQuery(req)
		response := generateResponse(s, searchQuery, w)

		movies_sort := movies.Movies(response["result"].([]movies.Movie))

		sort.Sort(movies.Movies(movies_sort))

		json.NewEncoder(w).Encode(movies_sort)

	}
}

func createSearchMoviesCompleteListHandler(s movies.MovieSearcher) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {

		page  := 1

		
		// get only one page
		searchQuery := generateQuery(req)
		response := generateResponse(s, searchQuery, w)
		r := response
		

		// get all pages
		for (len((response["result"].([]movies.Movie)))>0 ){

			if(page>1){
				r["result"] = append(r["result"].([]movies.Movie), response["result"].([]movies.Movie)...)
			}

			// Get next page
			page++
			q := req.URL.Query()
			if(q["p"]!=nil){
				q["p"][0] = strconv.Itoa(page)
			}else{
				q.Add("p", strconv.Itoa(page))
			}

		    req.URL.RawQuery = q.Encode() 
		    searchQuery = generateQuery(req)
			response = generateResponse(s, searchQuery, w)
		}

		json.NewEncoder(w).Encode(r)

	}
}

func generateQuery(req *http.Request) map[string]interface{}{

	// get parameter from request
	queryParams := req.URL.Query()
	searchQuery := make(map[string]interface{})
	searchQuery["q"] = queryParams.Get("q")
	if (queryParams.Get("p") != ""){
		searchQuery["p"] = queryParams.Get("p")
	}

	return searchQuery
}

func generateResponse(s movies.MovieSearcher, searchQuery map[string]interface{}, w http.ResponseWriter) map[string]interface{}{

	// call service to get movies
	movies, err := s.SearchMovies(searchQuery)

	// generate response
	response := make(map[string]interface{})
	response["result"] = movies
	if err != nil {
		response["error"] = err.Error()
	}

	return response
	
}

// NewHandler returns a router with all endpoint handlers
func NewHandler(s movies.MovieSearcher) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/movies", createSearchMoviesHandler(s)).Methods("GET")
	router.HandleFunc("/movies-sorted", createSearchMoviesSortHandler(s)).Methods("GET")
	router.HandleFunc("/movies-complete-list", createSearchMoviesCompleteListHandler(s)).Methods("GET")
	return router
}


