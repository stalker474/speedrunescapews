package webapp

//JSONUser blablabla
type JSONUser struct {
	Username string `json:"username"`
}

//JSONAuth blablabla
type JSONAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//JSONGameVictoryCondition blablabla
//Type can be QUEST or XP or LVL
//Name is the name of the skill or quest
//Value is the minimum value to achieve (must be "COMPLETED" for quest)
type JSONGameVictoryCondition struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

//JSONCreateChallenge blablabla
type JSONCreateChallenge struct {
	Username   string   `json:"username"`
	Opponent   string   `json:"opponent"`
	Game       string   `json:"game"`
	Challenges []string `json:"challenges"`
}

//JSONResult blablabla
type JSONResult struct {
	Result string `json:"result"`
}

//JSONUsersList list of users
type JSONUsersList struct {
	Users []string `json:"users"`
}

//JSONChallengesList blablabla
type JSONChallengesList struct {
	Result     string              `json:"result"`
	Challenges []JSONChallengeData `json:"challenges"`
}

//JSONSupportedGamesList blablabla
type JSONSupportedGamesList struct {
	Result string              `json:"result"`
	Games  []JSONSupportedGame `json:"games"`
}

//JSONSupportedGameChallenge blablabla
type JSONSupportedGameChallenge struct {
	Name    string `json:"name"`
	IconURL string `json:"iconurl"`
}

//JSONSupportedGame blablabla
type JSONSupportedGame struct {
	Name                string                       `json:"name"`
	SupportsHiscores    bool                         `json:"supportshiscores"`
	SupportsPlayTime    bool                         `json:"supportsplaytime"`
	SupportsAntiPay2Win bool                         `json:"supportsantipay2win"`
	SupportsAntiCoop    bool                         `json:"supportsanticoop"`
	SupportsQuests      bool                         `json:"supportsquests"`
	SupportsAutoAccount bool                         `json:"supportsautoaccount"`
	SupportsAntiCheat   bool                         `json:"supportsanticheat"`
	Website             string                       `json:"website"`
	PlayLink            string                       `json:"playlink"`
	IconURL             string                       `json:"iconurl"`
	Challenges          []JSONSupportedGameChallenge `json:"challenges"`
}

//JSONChallengeData blablabla
type JSONChallengeData struct {
	Status        string                   `json:"status"`
	Creator       string                   `json:"creator"`
	WinnerCreator bool                     `json:"winnercreator"`
	Game          JSONSupportedGame        `json:"game"`
	Name          string                   `json:"name"`
	ID            int                      `json:"id"`
	Private       JSONPrivateChallengeData `json:"private"`
}

//JSONPrivateChallengeData blablabla
type JSONPrivateChallengeData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

//JSONAccept blablabla
type JSONAccept struct {
	Username string `json:"username"`
	ID       int    `json:"id"`
}
