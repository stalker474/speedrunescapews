package main

// AccountPrivateData private data
type AccountPrivateData struct {
	Email    string
	Password string
}

// RSGameAccount Runescape account
type RSGameAccount struct {
	IsRS3    bool
	Username string
	Private  AccountPrivateData

	LiveData *HiscoresData
}

// CreateGameAccount Creates a new game account
func CreateGameAccount(IsRS3 bool) *RSGameAccount {
	acc := new(RSGameAccount)
	acc.IsRS3 = IsRS3
	//generate random name
	acc.Username = GenerateRandomName()
	//and password
	acc.Private.Password = GenerateRandomPassword()
	//and mail
	acc.Private.Email = GenerateMail(acc.Username)
	return acc
}

// LoadGameAccount Load a game account
func LoadGameAccount(IsRS3 bool, username string) *RSGameAccount {
	acc := new(RSGameAccount)
	acc.IsRS3 = IsRS3
	acc.Username = username
	return acc
}

func (acc *RSGameAccount) validateUserName() (bool, error) {
	return true, nil
}

// UpdateStats Update account statistics
func (acc *RSGameAccount) UpdateStats() error {
	fetcher := Hiscores{}
	var err error
	acc.LiveData, err = fetcher.Retrieve(acc.Username, HARDCORE, acc.IsRS3)
	if err != nil {
		return err
	}

	return nil
}
