package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/auth0-community/go-auth0"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	jose "gopkg.in/square/go-jose.v2"
)

const defaultPort = "8080"

//WSApp global application
var WSApp = App{}

func allowAll(origin string) bool {
	return true
}

/* Set up a global string for our secret */
var mySigningKey = []byte("secret")

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := []byte("YVpSiDm6qS7tT0cy2Lk_p8pSdbMMi8njVihwM0EDVd2e1ompEzxVETiT8b42b5m_")
		secretProvider := auth0.NewKeyProvider(secret)
		audience := []string{"https://speedrunescape.eu.auth0.com/userinfo"}

		configuration := auth0.NewConfiguration(secretProvider, audience, "https://speedrunescape.eu.auth0.com/", jose.HS256)
		validator := auth0.NewValidator(configuration, nil)

		token, err := validator.ValidateRequest(r)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Token is not valid:", token)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// StartWS Create a new game with a friend
func (*App) StartWS() {
	router := mux.NewRouter()
	router.Handle("/createuser", CreateUserHandler).Methods("POST")
	router.Handle("/deleteuser", authMiddleware(DeleteUserHandler)).Methods("POST")
	router.Handle("/authenticate", LoginUserHandler).Methods("POST")
	router.Handle("/finduser", GetUserByNameHandler).Methods("POST")
	router.Handle("/newchallenge", authMiddleware(CreateChallengeHandler)).Methods("POST")
	router.Handle("/mychallenges", authMiddleware(GetChallengesHandler)).Methods("POST")
	router.Handle("/accept", authMiddleware(AcceptHandler)).Methods("POST")
	router.Handle("/decline", authMiddleware(DeclineHandler)).Methods("POST")
	router.Handle("/terminate", authMiddleware(TerminateHandler)).Methods("POST")
	router.Handle("/supportedgames", GetSupportedGamesHandler).Methods("GET")

	c := cors.New(cors.Options{
		AllowOriginFunc:  allowAll,
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	logHandler := handlers.LoggingHandler(os.Stdout, handler)

	var port = defaultPort
	var externalPort = os.Getenv("PORT")
	if externalPort != "" {
		port = externalPort
	}

	http.ListenAndServe(":"+port, logHandler)
}

func respondJSONError(w http.ResponseWriter, msg string) {
	var res = JSONResult{}
	w.WriteHeader(http.StatusOK)
	res.Result = msg
	json.NewEncoder(w).Encode(res)
}

// CreateUserHandler WS
var CreateUserHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
})

// DeleteUserHandler WS
var DeleteUserHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
})

// LoginUserHandler WS
var LoginUserHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
})

// CreateChallengeHandler WS
var CreateChallengeHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
})

// GetUserByNameHandler WS
var GetUserByNameHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
})

// GetSupportedGamesHandler WS
var GetSupportedGamesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	var res = JSONSupportedGamesList{}
	res.Result = "OK"
	for _, g := range SupportedGames {
		res.Games = append(res.Games, g)
	}
	json.NewEncoder(w).Encode(res)
})

// GetChallengesHandler WS
var GetChallengesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
})

// AcceptHandler WS
var AcceptHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
})

// DeclineHandler WS
var DeclineHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
})

// TerminateHandler WS
var TerminateHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
})
