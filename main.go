package main

import (
	"log"
	"net/http"
	"zpl/handlers"
)

func main() {
	http.HandleFunc("/players", handlers.GetPostPlayers)
	http.HandleFunc("/players/", handlers.GetPutPlayer)
	http.HandleFunc("/teams", handlers.GetPostTeams)
	http.HandleFunc("/teams/", handlers.GetPutTeam)
	log.Println("Server running on port 5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
