package main

import (
	"log"
	"net/http"
	"zpl/handlers"
)

func main() {
	http.HandleFunc("players", handlers.GetPostPlayers)
	http.HandleFunc("players/{playerId}", handlers.GetPutPlayer)
	http.HandleFunc("players", handlers.GetPostTeams)
	http.HandleFunc("players/{playerId}", handlers.GetPutTeam)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
