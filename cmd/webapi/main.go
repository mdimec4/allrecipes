package main

import (
	"allrecipes.com_parser/pkg/allrecipes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

/*
REST API guides
- When should we use PUT and when should we use POST?
  http://restcookbook.com/HTTP%20Methods/put-vs-post/
- HTTP Status Codes
  http://www.restapitutorial.com/httpstatuscodes.html
- Rest api in GO sample
  https://www.thepolyglotdeveloper.com/2016/07/create-a-simple-restful-api-with-golang/
*/

func getRecipe(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	recipe, err := allrecipes.GetRecipe(params["id"])
	if err != nil {
		if strings.Contains(err.Error(), "GetRecipeInfo") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(recipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Println("str ", string(b))
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func getEnvConf(envName, defVal string) string {
	if v := os.Getenv(envName); v != "" {
		return v
		fmt.Println(envName, "=", v)
	}
	fmt.Println(envName, "=", defVal)
	return defVal
}

func main() {
	m := http.NewServeMux()
	router := mux.NewRouter()
	router.HandleFunc("/api/recipe/{id}", getRecipe).Methods("GET")

	m.Handle("/api/", router)
	fmt.Fprintf(os.Stderr, "%v\n", http.ListenAndServe(getEnvConf("ALRECIPE_PARSER_LISTEN_ADDR", ":4007"), m))
}
