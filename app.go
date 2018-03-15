package main

import (
	"strconv"

	// "log"
	"errors"

	"math"
)

// App main application
type App struct {
	// Db access to the data
	Db *Database
}

// Init construct the app object
func (a *App) Init() {
	a.Db = NewDatabase()
	a.Db.Connect()
	defer a.Db.Close()
	a.Db.Init()
}

// Close destroys the app object
func (a *App) Close() {

}

func (a *App) addUser(username string, password string) error {
	a.Db.Connect()
	defer a.Db.Close()
	if len(username) < 4 {
		return errors.New("User must be at least 4 characters long")
	}
	if err := a.findUser(username); err == nil {
		return errors.New("User already exists")
	}
	a.Db.AddUser(&User{username, password})
	return nil
}

func (a *App) removeUser(username string) error {
	a.Db.Connect()
	defer a.Db.Close()
	if user, err := a.Db.FindUser(username); err == nil {
		a.Db.RemoveUser(&user)
		return nil
	}
	return errors.New("User not found")
}

func (a *App) findUser(username string) error {
	a.Db.Connect()
	defer a.Db.Close()
	if _, err := a.Db.FindUser(username); err == nil {
		return nil
	}
	return errors.New("User not found")
}

func (a *App) authUser(username string, password string) error {
	a.Db.Connect()
	defer a.Db.Close()
	if user, err := a.Db.FindUser(username); err == nil {
		if user.Password == password {
			return nil
		}
		return errors.New("Invalid password")
	}
	return errors.New("User not found")
}

func (a *App) createChallenge(name string, username string, opponent string, gameT GameType) (int, error) {
	a.Db.Connect()
	defer a.Db.Close()
	if name == "" {
		return 0, errors.New("A challenge must have a name")
	}

	if username == opponent {
		return 0, errors.New("You cant challenge yourself")
	}

	if err := a.findUser(username); err != nil {
		return 0, errors.New("Username not found")
	}

	if err := a.findUser(opponent); err != nil {
		return 0, errors.New("Opponent not found")
	}

	challenge := Challenge{}
	challenge.ID = GenerateID()
	challenge.Name = name
	challenge.Completed = false

	challenge.Creator = username
	challenge.Opponent = opponent
	challenge.GameType = gameT
	challenge.WinnerCreator = false

	acc1 := CreateGameAccount(gameT == TYPERS3)
	gameAcc1 := GameAccount{acc1.Username, acc1.Private.Email, acc1.Private.Password}

	if err := a.Db.AddGameAccount(&gameAcc1); err != nil {
		return 0, errors.New("Failed to create game account")
	}

	acc2 := CreateGameAccount(gameT == TYPERS3)
	gameAcc2 := GameAccount{acc2.Username, acc2.Private.Email, acc2.Private.Password}

	if err := a.Db.AddGameAccount(&gameAcc2); err != nil {
		return 0, errors.New("Failed to create game account")
	}

	challenge.CreatorAccount = gameAcc1.Username
	challenge.OpponentAccount = gameAcc2.Username

	if err := a.Db.AddChallenge(&challenge); err != nil {
		//challenge creation failed, to avoid this "problem" use a transaction style request in the future
		a.Db.RemoveGameAccount(&gameAcc1)
		a.Db.RemoveGameAccount(&gameAcc2)
	}

	//challenge succesfully created and persisted
	return challenge.ID, nil
}

func (a *App) validateChallenge(id int, username string) error {
	a.Db.Connect()
	defer a.Db.Close()
	var challenge Challenge
	var g1, g2 GameAccount
	var gacc *RSGameAccount
	var err error
	if challenge, err = a.Db.FindChallenge(id); err != nil {
		return errors.New("Challenge not found")
	}
	if err := a.findUser(username); err != nil {
		return errors.New("Username not found")
	}
	if g1, err = a.Db.FindGameAccount(challenge.CreatorAccount); err != nil {
		return errors.New("Failed to find creators game account")
	}
	if g2, err = a.Db.FindGameAccount(challenge.OpponentAccount); err != nil {
		return errors.New("Failed to find opponents game account")
	}

	creatorStats := LoadGameAccount(challenge.GameType == TYPERS3, g1.Username)
	opponentStats := LoadGameAccount(challenge.GameType == TYPERS3, g2.Username)
	if creatorStats == nil || opponentStats == nil {
		return errors.New("Failed to init game accounts stats fetchers")
	}

	if err = creatorStats.UpdateStats(); err != nil {
		return errors.New("Failed to fetch opponents accounts stats. The account must be created and Hardcore Ironman mode selected")
	}

	if err = opponentStats.UpdateStats(); err != nil {
		return errors.New("Failed to fetch opponents accounts stats. The account must be created and Hardcore Ironman mode selected")
	}

	if username == challenge.Creator {
		gacc = creatorStats
	} else if username == challenge.Opponent {
		gacc = opponentStats
	} else {
		return errors.New("This user is neither the creator nor opponent of this challenge")
	}

	if challenge.Name == "Total Lvl" {
		var el SkillData
		var found bool
		if el, found = gacc.LiveData.Data["Overall"]; !found {
			return errors.New("Unexpected stats chart")
		}
		if el.Level < 800 {
			return errors.New("Failed to achieve total lvl 800, current lvl " + strconv.Itoa(el.Level))
		}
	} else if challenge.Name == "Combat Lvl" {
		var attack, defence, strength, hitpoints, prayer, ranged, magic SkillData
		var found bool
		if attack, found = gacc.LiveData.Data["Attack"]; !found {
			return errors.New("Unexpected stats chart")
		}
		if defence, found = gacc.LiveData.Data["Defence"]; !found {
			return errors.New("Unexpected stats chart")
		}
		if strength, found = gacc.LiveData.Data["Strength"]; !found {
			return errors.New("Unexpected stats chart")
		}
		if hitpoints, found = gacc.LiveData.Data["Hitpoints"]; !found {
			return errors.New("Unexpected stats chart")
		}
		if prayer, found = gacc.LiveData.Data["Prayer"]; !found {
			return errors.New("Unexpected stats chart")
		}
		if ranged, found = gacc.LiveData.Data["Ranged"]; !found {
			return errors.New("Unexpected stats chart")
		}
		if magic, found = gacc.LiveData.Data["Magic"]; !found {
			return errors.New("Unexpected stats chart")
		}

		var base = (float64(defence.Level+hitpoints.Level) + math.Floor(float64(prayer.Level)/2.0)) * 0.25
		var melee = float64(attack.Level+strength.Level) * 0.325
		var rang = math.Floor(float64(ranged.Level)*1.5) * 0.325
		var mage = math.Floor(float64(magic.Level)*1.5) * 0.325
		var max = math.Max(math.Max(melee, rang), mage)

		var level = int(math.Floor(((base+max)*100)+0.5) / 100)

		if level < 45 {
			return errors.New("Failed to achieve combat lvl 45, current lvl " + strconv.Itoa(level))
		}
	} else if challenge.Name == "Complete Demon Slayer quest" {
		if challenge.GameType != TYPERS3 {
			return errors.New("This challenge isnt available for this game")
		}
		found := false
		for _, a := range gacc.LiveData.CompletedQuests {
			if a == "Demon Slayer" {
				found = true
			}
		}
		if !found {
			return errors.New("The quest Demon Slayer isn't completed yet")
		}
	} else if challenge.Name == "Complete Dragon Slayer quest" {
		if challenge.GameType != TYPERS3 {
			return errors.New("This challenge isnt available for this game")
		}
		found := false
		for _, a := range gacc.LiveData.CompletedQuests {
			if a == "Dragon Slayer" {
				found = true
			}
		}
		if !found {
			return errors.New("The quest Dragon Slayer isn't completed yet")
		}
	} else {
		return errors.New("Fatal : Unknown challenge")
	}

	return nil
}
