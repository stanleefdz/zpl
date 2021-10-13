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

func handleSuccess() []byte {
	succModel := models.SuccessMessage{
		Status:  "Success",
		Message: "Successfully executed the API request",
	}
	b, _ := json.Marshal(succModel)
	return b
}

func writeResponse(rw http.ResponseWriter, code int, b []byte) {
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(code)
	rw.Write(b)
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

func GetTeamById(id int64) (*models.Team, error) {
	var team *models.Team
	res := teamCollection.FindOne(
		context.Background(),
		bson.D{{"_id", id}},
	)
	if res.Err() != nil {
		log.Println("Corresponding player not found")
		return nil, res.Err()
	}
	res.Decode(&team)
	return team, nil
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
			log.Printf("Players not found in db")
			writeResponse(rw, http.StatusBadRequest, handleError(fmt.Errorf("Players not found")))
			return
		}
		respByte, _ := json.Marshal(players)
		writeResponse(rw, http.StatusOK, respByte)
	} else if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error while fetching body, Error is: %s", err.Error())
			writeResponse(rw, http.StatusInternalServerError, handleError(err))
			return
		}
		var player models.Player
		err = json.Unmarshal(body, &player)
		if err != nil {
			log.Printf("Error while unmarshaling body, Error is: %s", err.Error())
			writeResponse(rw, http.StatusInternalServerError, handleError(err))
			return
		}
		p, _ := GetPlayerById(player.ID)
		if p != nil {
			log.Printf("Already exists collection with player Id")
			writeResponse(rw, http.StatusInternalServerError, handleError(fmt.Errorf("Player already exists")))
			return
		}
		result, err := playerCollection.InsertOne(
			context.Background(),
			bson.D{
				{"_id", player.ID},
				{"name", player.Name},
				{"age", player.Age},
				{"role", player.Role},
				{"country", player.Country},
				{"batting_style", player.BattingStyle},
				{"bowling_style", player.BowlingStyle},
			})
		if err != nil {
			log.Printf("Error while inserting collection to db")
			writeResponse(rw, http.StatusBadRequest, handleError(fmt.Errorf("Unable to insert collection to db")))
			return
		}
		log.Println(result.InsertedID)
		writeResponse(rw, http.StatusOK, handleSuccess())
	} else {
		writeResponse(rw, http.StatusMethodNotAllowed, handleError(fmt.Errorf("Method Not Allowed")))
	}
}

func GetPutPlayer(rw http.ResponseWriter, r *http.Request) {
	log.Println("Request Recieved")
	if r.Method == http.MethodGet {
		log.Println(r.URL.Path)
		ids := strings.TrimPrefix(r.URL.Path, "/players/")
		id, _ := strconv.ParseInt(ids, 10, 64)
		player, err := GetPlayerById(id)
		if err != nil {
			log.Printf("Player not found in db")
			writeResponse(rw, http.StatusBadRequest, handleError(fmt.Errorf("Player not found")))
			return
		}
		respByte, _ := json.Marshal(player)
		writeResponse(rw, http.StatusOK, respByte)
	} else if r.Method == http.MethodPut {
		log.Println(r.URL.Path)
		ids := strings.TrimPrefix(r.URL.Path, "/players/")
		id, _ := strconv.ParseInt(ids, 10, 64)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error while fetching body, Error is: %s", err.Error())
			writeResponse(rw, http.StatusInternalServerError, handleError(err))
			return
		}
		var player models.Player
		err = json.Unmarshal(body, &player)
		if err != nil {
			log.Printf("Error while unmarshaling body, Error is: %s", err.Error())
			writeResponse(rw, http.StatusInternalServerError, handleError(err))
			return
		}
		result, err := playerCollection.UpdateOne(
			context.Background(),
			bson.M{"_id": id},
			bson.D{
				{"$set", bson.D{
					{"name", player.Name},
					{"age", player.Age},
					{"role", player.Role},
					{"country", player.Country},
					{"batting_style", player.BattingStyle},
					{"bowling_style", player.BowlingStyle},
				}}},
		)
		if err != nil {
			log.Printf("Error while updating collection")
			writeResponse(rw, http.StatusBadRequest, handleError(err))
			return
		}
		log.Println(result.UpsertedID)
		writeResponse(rw, http.StatusOK, handleSuccess())
	} else {
		writeResponse(rw, http.StatusMethodNotAllowed, handleError(fmt.Errorf("Method Not Allowed")))
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
			log.Printf("Error while fetching teams from db, Error is: %s", err.Error())
			writeResponse(rw, http.StatusInternalServerError, handleError(err))
			return
		}
		for cursor.Next(ctx) {
			var team models.Team
			err = cursor.Decode(&team)
			if err != nil {
				log.Println(err)
			}
			teams = append(teams, team)
		}
		if len(teams) == 0 {
			log.Printf("Teams not found in db")
			writeResponse(rw, http.StatusBadRequest, handleError(fmt.Errorf("Teams not found")))
			return
		}
		respByte, _ := json.Marshal(teams)
		writeResponse(rw, http.StatusOK, respByte)
	} else if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error while fetching body, Error is: %s", err.Error())
			writeResponse(rw, http.StatusInternalServerError, handleError(err))
			return
		}
		var team models.Team
		err = json.Unmarshal(body, &team)
		if err != nil {
			log.Printf("Error while unmarshaling body, Error is: %s", err.Error())
			writeResponse(rw, http.StatusInternalServerError, handleError(err))
			return
		}
		t, _ := GetTeamById(team.ID)
		if t != nil {
			log.Printf("Already exists collection with team Id")
			writeResponse(rw, http.StatusInternalServerError, handleError(fmt.Errorf("Player already exists")))
			return
		}
		result, err := teamCollection.InsertOne(
			context.Background(),
			bson.D{
				{"_id", team.ID},
				{"name", team.Name},
				{"qwner", team.Owner},
				{"homeground", team.HomeGround},
			})
		if err != nil {
			log.Printf("Error while inserting collection to db")
			writeResponse(rw, http.StatusBadRequest, handleError(fmt.Errorf("Unable to insert collection to db")))
			return
		}
		log.Println(result.InsertedID)
		writeResponse(rw, http.StatusOK, handleSuccess())
	} else {
		writeResponse(rw, http.StatusMethodNotAllowed, handleError(fmt.Errorf("Method Not Allowed")))
	}
}

func GetPutTeam(rw http.ResponseWriter, r *http.Request) {
	log.Println("Request Recieved")
	if r.Method == http.MethodGet {
		log.Println(r.URL.Path)
		ids := strings.TrimPrefix(r.URL.Path, "/teams/")
		id, _ := strconv.ParseInt(ids, 10, 64)
		team, err := GetTeamById(id)
		if err != nil {
			log.Printf("Team not found in db")
			writeResponse(rw, http.StatusBadRequest, handleError(fmt.Errorf("Team not found")))
			return
		}
		respByte, _ := json.Marshal(team)
		writeResponse(rw, http.StatusOK, respByte)
	} else if r.Method == http.MethodPut {
		log.Println(r.URL.Path)
		ids := strings.TrimPrefix(r.URL.Path, "/teams/")
		id, _ := strconv.ParseInt(ids, 10, 64)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error while fetching body, Error is: %s", err.Error())
			writeResponse(rw, http.StatusInternalServerError, handleError(err))
			return
		}
		var team models.Team
		err = json.Unmarshal(body, &team)
		if err != nil {
			log.Printf("Error while unmarshaling body, Error is: %s", err.Error())
			writeResponse(rw, http.StatusInternalServerError, handleError(err))
			return
		}
		result, err := playerCollection.UpdateOne(
			context.Background(),
			bson.M{"_id": id},
			bson.D{
				{"$set", bson.D{
					{"name", team.Name},
					{"owner", team.Owner},
					{"homeground", team.HomeGround},
				}}},
		)
		if err != nil {
			log.Printf("Error while updating collection")
			writeResponse(rw, http.StatusBadRequest, handleError(err))
			return
		}
		log.Println(result.UpsertedID)
		writeResponse(rw, http.StatusOK, handleSuccess())
	} else {
		writeResponse(rw, http.StatusMethodNotAllowed, handleError(fmt.Errorf("Method Not Allowed")))
	}
}
