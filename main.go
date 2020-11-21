package main

import (
    "bytes"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
    "github.com/joho/godotenv"
    "time"
    "github.com/kalyani14377/weather-app/weather"
)


var tpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func main() {
    err := godotenv.Load()
    if err != nil {
		log.Println("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
    }
    //Weather application
    apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("Env: apiKey must be set")
	}

    myClient := &http.Client{Timeout: 10 * time.Second}
    weatherapi := weather.NewClient(myClient, apiKey)

    fs := http.FileServer(http.Dir("assets"))
    mux := http.NewServeMux()

	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
    // mux.HandleFunc("/", indexHandler)
    mux.HandleFunc("/", searchHandler(weatherapi))
	http.ListenAndServe(":"+port, mux)
}

func searchHandler(weatherapi *weather.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()
        searchQuery := params.Get("location")
        
        results, err := weatherapi.FetchWeather(searchQuery)
        log.Println(results);
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}


		search := &Search{
			Query:      searchQuery,
			Results:    results,
		}

		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf.WriteTo(w)
	}
}