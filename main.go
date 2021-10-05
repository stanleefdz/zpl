package main

import (
	"net/http"
	"zpl/handlers"
)

func main() {
	http.HandleFunc("players/{playerId}", handlers.GetPlayers)
}
