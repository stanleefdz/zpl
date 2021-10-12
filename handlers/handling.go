package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func handleError(err error) []byte {
	errModel := models.ErrorMessage{
		Description: err.Error(),
		Message:     "Invalid Request Params",
	}
	b, _ := json.Marshal(errModel)
	return b
}

func writeResponse(rw http.ResponseWriter, code int, b []byte) {
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(code)
	rw.Write(b)
}

func GetPostPlayers(rw http.ResponseWriter, r *http.Request) {
	log.Println("Request Recieved")
	if r.Method == http.MethodGet {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var players []models.Player
		cursor, err := playerCollection.Find(ctx, bson.D{})
		if err != nil {
			log.Printf("Error while fetching players from db, Error is: %s", err.Error())
			writeResponse(rw, http.StatusInternalServerError, handleError(err))
			return
		}
		for cursor.Next(ctx) {
			var player models.Player
			err = cursor.Decode(&player)
			if err != nil {
				log.Println(err)
			}
			players = append(players, player)
		}
		if len(players) == 0 {
			log.Printf("players not found in db")
			writeResponse(rw, http.StatusBadRequest, handleError(fmt.Errorf("players not found")))
			return
		}
		respByte, _ := json.Marshal(players)
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respByte)
	} else if r.Method == http.MethodPost {
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
				{"_id", player.ID},
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
		log.Println(result.InsertedID)
		resp := "Success"
		respByte, _ := json.Marshal(resp)
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respByte)
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write(handleError(fmt.Errorf("Method Not Allowed")))
	}
}

func GetPlayerById(id int64) (*models.Player, error) {
	var player *models.Player
	res := playerCollection.FindOne(
		context.Background(),
		bson.D{{"_id", id}},
	)
	if res.Err() != nil {
		log.Println("Corresponding player not found")
		return nil, res.Err()
	}
	res.Decode(&player)
	return player, nil
}

func GetPutPlayer(rw http.ResponseWriter, r *http.Request) {
	log.Println("Request Recieved")
	if r.Method == http.MethodGet {
		log.Println(r.URL.Path)
		ids := strings.TrimPrefix(r.URL.Path, "/players/")
		id, _ := strconv.ParseInt(ids, 10, 64)
		player, err := GetPlayerById(id)
		if err != nil {
			writeResponse(rw, http.StatusBadRequest, handleError(err))
			return
		}
		b, _ := json.Marshal(player)
		writeResponse(rw, 200, b)
	} else if r.Method == http.MethodPut {
		id := r.URL.Query().Get("playerId")
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
				{"ID", id},
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
		log.Println(result.InsertedID)
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

func GetPostTeams(rw http.ResponseWriter, r *http.Request) {
	log.Println("Request Recieved")
	if r.Method == http.MethodGet {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var teams []models.Team
		cursor, err := teamCollection.Find(ctx, bson.D{})
		if err != nil {
			log.Println("400 Error")
		}
		for cursor.Next(ctx) {
			var team models.Team
			err = cursor.Decode(team)
			if err != nil {
				log.Println("500 Error")
			}
			teams = append(teams, team)
		}
		respByte, _ := json.Marshal(teams)
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respByte)
		rw.WriteHeader(http.StatusOK)
	} else if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("400 Error")
		}
		var team models.Team
		err = json.Unmarshal(body, &team)
		if err != nil {
			log.Println("400 Error")
		}
		result, err := teamCollection.InsertOne(
			context.Background(),
			bson.D{
				{"ID", team.ID},
				{"Name", team.Name},
				{"Owner", team.Owner},
				{"HomeGround", team.HomeGround},
			})
		if err != nil {
			log.Println("500 Error")
		}
		log.Println(result.InsertedID)
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

func GetTeamById(id string) (models.Team, error) {
	var team models.Team
	cursor, err := teamCollection.Find(
		context.Background(),
		bson.D{{"ID", id}},
	)
	if err != nil {
		log.Println("Corresponding player not found")
		return team, err
	}
	cursor.Decode(&team)
	return team, nil
}

func GetPutTeam(rw http.ResponseWriter, r *http.Request) {
	log.Println("Request Recieved")
	if r.Method == http.MethodGet {
		id := r.URL.Query().Get("teamId")
		team, err := GetTeamById(id)
		if err != nil {
			rw.Write([]byte("Team Not Found"))
			rw.WriteHeader(http.StatusNoContent)
		}
		respByte, _ := json.Marshal(team)
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respByte)
		rw.WriteHeader(http.StatusOK)
	} else if r.Method == http.MethodPut {
		id := r.URL.Query().Get("teamId")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("400 Error")
		}
		var team models.Team
		err = json.Unmarshal(body, &team)
		if err != nil {
			log.Println("400 Error")
		}
		result, err := playerCollection.InsertOne(
			context.Background(),
			bson.D{
				{"ID", id},
				{"Name", team.Name},
				{"Owner", team.Owner},
				{"HomeGround", team.HomeGround},
			})
		if err != nil {
			log.Println("500 Error")
		}
		log.Println(result.InsertedID)
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
