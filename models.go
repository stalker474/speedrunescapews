package main

// User an app user
type User struct {
	Name     string
	Password string
}

// GameType available games
type GameType int

// GameType types
const (
	TYPEOSRS = GameType(0)
	TYPERS3  = GameType(1)
)

// GameState status of this game
type GameState int

// GameState types
const (
	WAITING   = GameState(0)
	STARTED   = GameState(1)
	COMPLETED = GameState(2)
)

// Challenge a challenge created
type Challenge struct {
	ID              int
	Creator         string
	Opponent        string
	Completed       bool
	WinnerCreator   bool
	GameType        GameType
	Name            string
	GameState       GameState
	CreatorAccount  string
	OpponentAccount string
}

type accountPrivateData struct {
}

// GameAccount an account ingame
type GameAccount struct {
	Username string
	Email    string
	Password string
}
