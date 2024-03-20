package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type Weather struct {
	Status Status `json:"status"`
}

func (w *Weather) checkStatus() (resWater string, resWind string) {
	switch {
	case w.Status.Water < 5:
		resWater = "Aman"
	case w.Status.Water >= 6 && w.Status.Water <= 8:
		resWater = "Siaga"
	case w.Status.Water > 8:
		resWater = "Bahaya"
	}

	switch {
	case w.Status.Wind < 6:
		resWind = "Aman"
	case w.Status.Wind >= 7 && w.Status.Wind <= 15:
		resWind = "Siaga"
	case w.Status.Wind > 15:
		resWind = "Bahaya"
	}
	return resWater, resWind
}

func generateJSON() {
	for {
		time.Sleep(15 * time.Second)

		weather := Weather{
			Status: Status{
				Water: rand.Intn(100) + 1, // Generate random value between 1-100 for water
				Wind:  rand.Intn(100) + 1, // Generate random value between 1-100 for wind
			},
		}

		file, err := os.Create("weather.json")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		err = encoder.Encode(weather)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("weather.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var weather Weather
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&weather)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resWater, resWind := weather.checkStatus()

	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Weather Status</title>
	</head>
	<body>
	<h1>Weather Status</h1>
	<p>Water: %d - %s</p>
	<p>Wind: %d - %s</p>
	<button onclick="refreshPage()">Refresh</button>

    <script>
        function refreshPage() {
            location.reload();
        }
    </script>

	</body>
	</html>`, weather.Status.Water, resWater, weather.Status.Wind, resWind)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, html)
}

func main() {
	go generateJSON() // Start goroutine to generate JSON

	http.HandleFunc("/", statusHandler)

	fmt.Println("Server is running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
