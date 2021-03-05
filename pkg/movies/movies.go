package movies

// Movie represents a single movie
type Movie struct {
	Title string `json:"Title"`
	Year  string `json:"Year"`
}

type Movies []Movie

func (m Movies) Len() int           { return len(m) }
func (m Movies) Less(i, j int) bool { 
	if(m[i].Year == m[j].Year){
		return m[i].Title < m[j].Title
	}else{
		return m[i].Year < m[j].Year 
	} 
}
func (m Movies) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

// MovieSearcher is the interfaces for anything that searches for movies
type MovieSearcher interface {
	SearchMovies(query map[string]interface{}) ([]Movie, error)
}
