package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"zpl/db"
	"zpl/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var dbclient *mongo.Client
var playerCollection *mongo.Collection
var teamCollection *mongo.Collection

func init() {
	dbclient = db.ConnectionHandler()
	playerCollection = dbclient.Database("zpl").Collection("players")
	teamCollection = dbclient.Database("zpl").Collection("teams")
}

func GetPostPlayers(rw http.ResponseWriter, r *http.Request) {
	log.Println("Request Recieved")
	if r.Method == http.MethodGet {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var players []models.Player
		cursor, err := playerCollection.Find(ctx, bson.D{})
		if err != nil {
			log.Println("400 Error")
		}
		for cursor.Next(ctx) {
			var player models.Player
			err = cursor.Decode(player)
			if err != nil {
				log.Println("500 Error")
			}
			players = append(players, player)
		}
		respByte, _ := json.Marshal(players)
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respByte)
		rw.WriteHeader(http.StatusOK)
	} else if r.Method == http.MethodPost {

	} else {
		rw.Write([]byte("Method Not Allowed"))
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func GetById(id string) (models.Player, error) {
	var player models.Player
	cursor, err := playerCollection.Find(
		context.Background(),
		bson.D{{"ID", id}},
	)
	if err != nil {
		log.Println("Corresponding player not found")
		return player, err
	}
	cursor.Decode(&player)
	return player, nil
}

func GetPutPlayer(rw http.ResponseWriter, r *http.Request) {
	log.Println("Request Recieved")
	if r.Method == http.MethodGet {
		id := r.URL.Query().Get("playerId")
		player, err := GetById(id)
		if err != nil {
			rw.Write([]byte("Player Not Found"))
			rw.WriteHeader(http.StatusNoContent)
		}
		respByte, _ := json.Marshal(player)
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respByte)
		rw.WriteHeader(http.StatusOK)
	} else if r.Method == http.MethodPut {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("400 Error")
		}
		var player models.Player
		err = json.Unmarshal(body, &player)
		if err != nil {
			log.Println("400 Error")
		}
		result, err := playerCollection.InsertOne(
			context.Background(),
			bson.D{
				{"ID", player.ID},
				{"Name", player.Name},
				{"Age", player.Age},
				{"Role", player.Role},
				{"Country", player.Country},
				{"BattingStyle", player.BattingStyle},
				{"BowlingStyle", player.BowlingStyle},
			})
		if err != nil {
			log.Println("500 Error")
		}
		resp := "Success"
		respByte, _ := json.Marshal(resp)
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respByte)
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.Write([]byte("Method Not Allowed"))
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}
