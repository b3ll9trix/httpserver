package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// Public Members
type Joke struct {
	Joke string `json:"joke"`
}

func GetRandomJoke(w http.ResponseWriter, r *http.Request) {
	log.Printf("Fetching Joke...")
	w.Write([]byte(getJoke()))
}

// Private Members
func getJoke() string {
	joke := Joke{Joke: "bad time for a joke"}
	url := "https://v2.jokeapi.dev/joke/Any?type=single"
	resp, err := http.Get(url)
	if err != nil {
		return joke.Joke
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&joke)

	return joke.Joke

}
