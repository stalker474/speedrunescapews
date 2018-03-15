package main

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/go-sql-driver/mysql" //required by sql
)

const driverType = "mysql"
const user = "u3uhtccyht5t4ysp:3qH70xpRHfiilzUqkLU@tcp(b0stcadoi-mysql.services.clever-cloud.com:3306)/"
const dbName = "b0stcadoi"

//Database our main database for storing accounts
type Database struct {
	db *sql.DB
}

//NewDatabase Create new Database
func NewDatabase() *Database {
	return &Database{nil}
}

// Init Create tables
func (thisDatabase *Database) Init() {
	if thisDatabase.db == nil {
		log.Fatal("Database is not connected")
	}
	result, err := thisDatabase.db.Exec(`CREATE TABLE IF NOT EXISTS Accounts (
		username VARCHAR(45) NOT NULL, 
		password VARCHAR(45) NOT NULL, 
		PRIMARY KEY(username))`)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)

	result, err = thisDatabase.db.Exec(`CREATE TABLE IF NOT EXISTS GameAccount (
		username VARCHAR(45) UNIQUE NOT NULL,
		email VARCHAR(45)  NOT NULL,
		password VARCHAR(45)  NOT NULL,
		PRIMARY KEY (username))`)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)

	result, err = thisDatabase.db.Exec(`CREATE TABLE IF NOT EXISTS Challenges (
		id INT NOT NULL AUTO_INCREMENT,
		name VARCHAR(45) NOT NULL,
		creator VARCHAR(45) NOT NULL,
		opponent VARCHAR(45) NOT NULL,
		winnercreator TINYINT NULL,
		gametype VARCHAR(45) NOT NULL,
		gamestate VARCHAR(45) NOT NULL,
		creatorAccount VARCHAR(45) NULL,
		opponentAccount VARCHAR(45) NULL,
		PRIMARY KEY (id),
		INDEX creatorAccount_idx (creatorAccount ASC),
		INDEX opponentAccount_idx (opponentAccount ASC),
		CONSTRAINT creatorAccount
		  FOREIGN KEY (creatorAccount)
		  REFERENCES GameAccount (username)
		  ON DELETE CASCADE
		  ON UPDATE CASCADE,
		CONSTRAINT opponentAccount
		  FOREIGN KEY (opponentAccount)
		  REFERENCES GameAccount (username)
		  ON DELETE CASCADE
		  ON UPDATE CASCADE)`)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)
}

// Reset Create tables
func (thisDatabase *Database) Reset() {
	if thisDatabase.db == nil {
		log.Fatal("Database is not connected")
	}
	thisDatabase.db.Query("DROP TABLE Accounts")
	thisDatabase.db.Query("DROP TABLE GameAccount")
	thisDatabase.db.Query("DROP TABLE Challenges")
	thisDatabase.Init()
}

// Connect Connect to the database
func (thisDatabase *Database) Connect() {
	db, err := sql.Open(driverType, user)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("USE " + dbName)
	if err != nil {
		log.Fatal(err)
	}
	thisDatabase.db = db
}

// Close Disconnect from the database
func (thisDatabase *Database) Close() {
	if thisDatabase.db != nil {
		thisDatabase.db.Close()
	}
}

func (thisDatabase *Database) checkConnection() {
	if thisDatabase.db == nil {
		log.Fatal("Database is not connected")
	} else {
		thisDatabase.db.Exec("USE " + dbName)
	}
}

func validateUser(user *User) {
	if user == nil {
		log.Fatal("User is nil!")
	}
	if len(user.Name) < 1 {
		log.Fatal("Username isnt long enough!")
	}
	if len(user.Password) < 1 {
		log.Fatal("Password isnt long enough!")
	}
}

func validateGameAccount(acc *GameAccount) {
	if acc == nil {
		log.Fatal("Account is nil!")
	}
	if len(acc.Username) < 1 {
		log.Fatal("Username isnt long enough!")
	}
	if len(acc.Password) < 1 {
		log.Fatal("Password isnt long enough!")
	}
	if len(acc.Email) < 1 {
		log.Fatal("Email isnt long enough!")
	}
}

func validateChallenge(ch *Challenge) {
	if ch == nil {
		log.Fatal("Challenge is nil!")
	}
	if len(ch.Name) < 1 {
		log.Fatal("Name isnt long enough!")
	}
	if len(ch.Creator) < 1 {
		log.Fatal("Creator name isnt long enough!")
	}
	if len(ch.Opponent) < 1 {
		log.Fatal("Opponent name isnt long enough!")
	}
}

// AddUser Add a new account to persist
func (thisDatabase *Database) AddUser(user *User) error {
	validateUser(user)
	thisDatabase.checkConnection()
	_, err := thisDatabase.db.Query("INSERT INTO Accounts (Username,Password) VALUES (?,?)", user.Name, user.Password)
	return err
}

// RemoveUser Delete an account from persist
func (thisDatabase *Database) RemoveUser(user *User) error {
	validateUser(user)
	thisDatabase.checkConnection()
	_, err := thisDatabase.db.Query("DELETE FROM Accounts WHERE Username = ?", user.Name)
	return err
}

// UpdateUser Update an account from persist
func (thisDatabase *Database) UpdateUser(user *User) error {
	validateUser(user)
	thisDatabase.checkConnection()
	_, err := thisDatabase.db.Query("UPDATE Accounts SET Password = ? WHERE Username = ?", user.Password, user.Name)
	return err
}

// FindUser Find an account by name
func (thisDatabase *Database) FindUser(Name string) (User, error) {
	thisDatabase.checkConnection()
	rows, err := thisDatabase.db.Query("SELECT Username, Password FROM Accounts WHERE Username = ?", Name)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var username string
	var password string

	for rows.Next() {

		err2 := rows.Scan(&username, &password)
		if err2 != nil {
			log.Fatal(err2)
		}
		return User{username, password}, nil
	}
	return User{}, errors.New("User does not exist")
}

// AddGameAccount Add a new game account to persist
func (thisDatabase *Database) AddGameAccount(acc *GameAccount) error {
	validateGameAccount(acc)
	thisDatabase.checkConnection()
	_, err := thisDatabase.db.Query("INSERT INTO GameAccount (username,email,password) VALUES (?,?,?)", acc.Username, acc.Email, acc.Password)
	return err
}

// RemoveGameAccount Delete a game account from persist
func (thisDatabase *Database) RemoveGameAccount(acc *GameAccount) error {
	validateGameAccount(acc)
	thisDatabase.checkConnection()
	_, err := thisDatabase.db.Query("DELETE FROM GameAccount WHERE username = ?", acc.Username)
	return err
}

// FindGameAccount Find a game account by username
func (thisDatabase *Database) FindGameAccount(username string) (GameAccount, error) {
	thisDatabase.checkConnection()
	rows, err := thisDatabase.db.Query("SELECT username, email, password FROM GameAccount WHERE username = ?", username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	acc := GameAccount{}

	for rows.Next() {

		err2 := rows.Scan(&acc.Username, &acc.Email, &acc.Password)
		if err2 != nil {
			log.Fatal(err2)
		}
		return acc, nil
	}
	return GameAccount{}, errors.New("User does not exist")
}

// AddChallenge Add a new game account to persist
func (thisDatabase *Database) AddChallenge(cha *Challenge) error {
	validateChallenge(cha)
	thisDatabase.checkConnection()
	_, err := thisDatabase.db.Query("INSERT INTO Challenges (id, name ,creator, opponent, winnercreator, gameType, gameState, creatorAccount, opponentAccount) VALUES (?,?,?,?,?,?,?,?,?)",
		cha.ID, cha.Name, cha.Creator, cha.Opponent, cha.WinnerCreator, cha.GameType, cha.GameState, cha.CreatorAccount, cha.OpponentAccount)
	return err
}

// RemoveChallenge Delete a game account from persist
func (thisDatabase *Database) RemoveChallenge(cha *Challenge) error {
	validateChallenge(cha)
	thisDatabase.checkConnection()
	_, err := thisDatabase.db.Query("DELETE FROM Challenges WHERE id = ?", cha.ID)
	return err
}

// UpdateChallenge Update a challenge from persist
func (thisDatabase *Database) UpdateChallenge(cha *Challenge) error {
	validateChallenge(cha)
	thisDatabase.checkConnection()
	_, err := thisDatabase.db.Query("UPDATE Challenges SET winnercreator = ?, gameState = ? WHERE id = ?", cha.WinnerCreator, cha.GameState, cha.ID)
	return err
}

// FindChallenge Find a challenge by id
func (thisDatabase *Database) FindChallenge(id int) (Challenge, error) {
	thisDatabase.checkConnection()
	rows, err := thisDatabase.db.Query("SELECT id, name ,creator, opponent, winnercreator, gameType, gameState, creatorAccount, opponentAccount FROM Challenges WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	cha := Challenge{}

	for rows.Next() {

		err2 := rows.Scan(&cha.ID, &cha.Name, &cha.Creator, &cha.Opponent, &cha.WinnerCreator, &cha.GameType, &cha.GameState, &cha.CreatorAccount, &cha.OpponentAccount)
		if err2 != nil {
			log.Fatal(err2)
		}
		return cha, nil
	}
	return Challenge{}, errors.New("Challenge does not exist")
}

// FindChallengesByUser Find a challenge by creator or opponent username
func (thisDatabase *Database) FindChallengesByUser(username string) ([]Challenge, error) {
	thisDatabase.checkConnection()
	rows, err := thisDatabase.db.Query("SELECT id, name ,creator, opponent, winnercreator, gameType, gameState, creatorAccount, opponentAccount FROM Challenges WHERE creator = ? OR opponent = ?", username, username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	list := []Challenge{}
	cha := Challenge{}

	for rows.Next() {

		err2 := rows.Scan(&cha.ID, &cha.Name, &cha.Creator, &cha.Opponent, &cha.WinnerCreator, &cha.GameType, &cha.GameState, &cha.CreatorAccount, &cha.OpponentAccount)
		if err2 != nil {
			log.Fatal(err2)
		}
		list = append(list, cha)
	}
	if len(list) > 0 {
		return list, nil
	}
	return nil, errors.New("No challenges yet")
}
