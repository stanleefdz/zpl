package main

import (
	"net/http"
	"zpl/handlers"
)

func main() {
	http.HandleFunc("players", handlers.GetPostPlayers)
	//http.HandleFunc("players/{playerId}", handlers.GetPutPlayer)
}
