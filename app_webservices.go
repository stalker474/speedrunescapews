package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

//WSApp global application
var WSApp = App{}

func allowAll(origin string) bool {
	return true
}

// StartWS Create a new game with a friend
func (*App) StartWS() {
	router := mux.NewRouter()
	router.HandleFunc("/createuser", CreateUser).Methods("POST")
	router.HandleFunc("/deleteuser", DeleteUser).Methods("POST")
	router.HandleFunc("/authenticate", LoginUser).Methods("POST")
	router.HandleFunc("/finduser", GetUserByName).Methods("POST")
	router.HandleFunc("/newchallenge", CreateChallenge).Methods("POST")
	router.HandleFunc("/mychallenges", GetChallenges).Methods("POST")
	router.HandleFunc("/accept", Accept).Methods("POST")
	router.HandleFunc("/decline", Decline).Methods("POST")
	router.HandleFunc("/terminate", Terminate).Methods("POST")
	router.HandleFunc("/supportedgames", GetSupportedGames).Methods("GET")

	c := cors.New(cors.Options{
		AllowOriginFunc:  allowAll,
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	http.ListenAndServe(":8080", handler)
}

func respondJSONError(w http.ResponseWriter, msg string) {
	var res = JSONResult{}
	w.WriteHeader(http.StatusOK)
	res.Result = msg
	json.NewEncoder(w).Encode(res)
}

// CreateUser WS
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var usr JSONAuth
	var res = JSONResult{"OK"}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&usr); err != nil {
		respondJSONError(w, err.Error())
		return
	}

	if err := WSApp.addUser(usr.Username, usr.Password); err != nil {
		respondJSONError(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// DeleteUser WS
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var usr JSONUser
	var res = JSONResult{"OK"}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&usr); err != nil {
		respondJSONError(w, err.Error())
		return
	}

	if err := WSApp.removeUser(usr.Username); err != nil {
		respondJSONError(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// LoginUser WS
func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var usr JSONAuth
	var res = JSONResult{"OK"}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&usr); err != nil {
		respondJSONError(w, err.Error())
		return
	}

	if err := WSApp.authUser(usr.Username, usr.Password); err != nil {
		respondJSONError(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// CreateChallenge WS
func CreateChallenge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var createGame JSONCreateChallenge
	var res = JSONResult{"OK"}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&createGame); err != nil {
		respondJSONError(w, err.Error())
		return
	}

	var gType GameType
	switch createGame.Game {
	case "Old School Runescape":
		gType = TYPEOSRS
	case "Runescape3":
		gType = TYPERS3
	default:
		respondJSONError(w, "Game not supported")
		return
	}
	_, err := WSApp.createChallenge(createGame.Challenges[0], createGame.Username, createGame.Opponent, gType)
	if err != nil {
		respondJSONError(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// GetUserByName WS
func GetUserByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var usr JSONUser
	var res = JSONResult{"OK"}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&usr); err != nil {
		respondJSONError(w, err.Error())
		return
	}

	if err := WSApp.findUser(usr.Username); err != nil {
		respondJSONError(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// GetSupportedGames WS
func GetSupportedGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	var res = JSONSupportedGamesList{}
	res.Result = "OK"
	for _, g := range SupportedGames {
		res.Games = append(res.Games, g)
	}
	json.NewEncoder(w).Encode(res)
}

// GetChallenges WS
func GetChallenges(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var usr JSONUser
	var res JSONChallengesList
	res.Result = "OK"

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&usr); err != nil {
		respondJSONError(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	challenges, err := WSApp.Db.FindChallengesByUser(usr.Username)
	if err == nil {
		for _, challenge := range challenges {
			data := JSONChallengeData{}
			data.ID = challenge.ID
			switch challenge.GameState {
			case WAITING:
				data.Status = "WAITING"
			case STARTED:
				data.Status = "STARTED"
			case COMPLETED:
				data.Status = "COMPLETED"
			default:
				data.Status = "UNKNOWN"
			}
			data.Creator = challenge.Creator
			data.Name = challenge.Name
			data.WinnerCreator = challenge.WinnerCreator

			switch challenge.GameType {
			case TYPERS3:
				data.Game = SupportedGames[0]
			case TYPEOSRS:
				data.Game = SupportedGames[1]
			}

			if challenge.GameState >= STARTED {
				var acc GameAccount
				var erro error
				if challenge.Creator == usr.Username {
					acc, erro = WSApp.Db.FindGameAccount(challenge.CreatorAccount)
				} else {
					acc, erro = WSApp.Db.FindGameAccount(challenge.OpponentAccount)
				}
				if erro != nil {
					respondJSONError(w, err.Error())
					return
				}
				data.Private.Email = acc.Email
				data.Private.Password = acc.Password
				data.Private.Username = acc.Username
			}
			res.Challenges = append(res.Challenges, data)
		}
	} else {
		respondJSONError(w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(res)
}

// Accept WS
func Accept(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var usr JSONAccept
	var res JSONResult
	res.Result = "OK"

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&usr); err != nil {
		respondJSONError(w, err.Error())
		return
	}

	if challenge, err := WSApp.Db.FindChallenge(usr.ID); err == nil {
		if challenge.Opponent != usr.Username {
			respondJSONError(w, "Only the opponent can accept")
			return
		}
		if challenge.GameState != WAITING {
			respondJSONError(w, "This game is not in waiting mode")
			return
		}
		w.WriteHeader(http.StatusOK)
		//change game state
		challenge.GameState = STARTED
		if err := WSApp.Db.UpdateChallenge(&challenge); err != nil {
			respondJSONError(w, err.Error())
			return
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	respondJSONError(w, "Game not found")
}

// Decline WS
func Decline(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var usr JSONAccept
	var res JSONResult
	res.Result = "OK"

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&usr); err != nil {
		respondJSONError(w, err.Error())
		return
	}

	if challenge, err := WSApp.Db.FindChallenge(usr.ID); err == nil {
		if challenge.Opponent != usr.Username && challenge.Creator != usr.Username {
			respondJSONError(w, "Only a participant can decline")
			return
		}
		if challenge.GameState != WAITING {
			respondJSONError(w, "This game is not in waiting mode")
			return
		}
		w.WriteHeader(http.StatusOK)
		//change game state
		challenge.GameState = STARTED

		if err := WSApp.Db.RemoveChallenge(&challenge); err != nil {
			respondJSONError(w, err.Error())
			return
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	respondJSONError(w, "Game not found")
}

// Terminate WS
func Terminate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var usr JSONAccept
	var res JSONResult
	res.Result = "OK"

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&usr); err != nil {
		respondJSONError(w, err.Error())
		return
	}

	if challenge, err := WSApp.Db.FindChallenge(usr.ID); err == nil {
		if challenge.GameState != STARTED {
			respondJSONError(w, "This game is not started yet")
			return
		}
		if challenge.GameState == COMPLETED {
			respondJSONError(w, "This game is already completed")
			return
		}
		if err = WSApp.validateChallenge(challenge.ID, usr.Username); err != nil {
			respondJSONError(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		//change game state
		challenge.GameState = COMPLETED
		challenge.WinnerCreator = usr.Username == challenge.Creator
		if err := WSApp.Db.UpdateChallenge(&challenge); err != nil {
			respondJSONError(w, err.Error())
			return
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	respondJSONError(w, "Game not found")
}
