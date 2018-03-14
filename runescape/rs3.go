package runescape

import (
	"encoding/json"
	"log"
	"strings"

	"golang.org/x/net/html"
)

//RS3 data fetcher for runescape
type RS3 struct {
}

// RS3JsonData json data
type RS3JsonData struct {
	Magic            int `json:"magic"`
	Questsstarted    int `json:"questsstarted"`
	Totalskill       int `json:"totalskill"`
	Questscomplete   int `json:"questscomplete"`
	Questsnotstarted int `json:"questsnotstarted"`
	Totalxp          int `json:"totalxp"`
	Ranged           int `json:"ranged"`
	Activities       []struct {
		Date    string `json:"date"`
		Details string `json:"details"`
		Text    string `json:"text"`
	} `json:"activities"`
	Skillvalues []struct {
		Level int `json:"level"`
		Xp    int `json:"xp"`
		Rank  int `json:"rank"`
		ID    int `json:"id"`
	} `json:"skillvalues"`
	Name        string `json:"name"`
	Rank        string `json:"rank"`
	Melee       int    `json:"melee"`
	Combatlevel int    `json:"combatlevel"`
	LoggedIn    string `json:"loggedIn"`
}

// RS3JsonQuestsData json quests data
type RS3JsonQuestsData struct {
	Quests []struct {
		Title        string `json:"title"`
		Status       string `json:"status"`
		Difficulty   int    `json:"difficulty"`
		Members      bool   `json:"members"`
		QuestPoints  int    `json:"questPoints"`
		UserEligible bool   `json:"userEligible"`
	} `json:"quests"`
	LoggedIn string `json:"loggedIn"`
}

var skillsIDMap = map[int]string{
	15: "Herblore",
	2:  "Strength",
	3:  "Constitution",
	0:  "Attack",
	1:  "Defence",
	4:  "Ranged",
	18: "Slayer",
	6:  "Magic",
	23: "Summoning",
	10: "Fishing",
	5:  "Prayer",
	24: "Dungeoneering",
	13: "Smithing",
	17: "Thieving",
	16: "Agility",
	7:  "Cooking",
	11: "Firemaking",
	9:  "Flething",
	12: "Crafting",
	8:  "Woodcutting",
	14: "Mining",
	22: "Construction",
	20: "Runecrafting",
	19: "Farming",
	21: "Hunter",
	25: "Divination",
	26: "Invention",
}

// Retrieve Simply fetch the skill and completed quests data from json string
func (*RS3) Retrieve(jsonData string, jsonQuestsData string, hardcoreIronmanWebpage string) ([]SkillData, []string, bool, error) {
	var skilldata RS3JsonData
	var questsdata RS3JsonQuestsData
	dec := json.NewDecoder(strings.NewReader(jsonData))
	err := dec.Decode(&skilldata)
	if err != nil {
		return nil, nil, false, err
	}

	dec = json.NewDecoder(strings.NewReader(jsonQuestsData))
	err = dec.Decode(&questsdata)
	if err != nil {
		return nil, nil, false, err
	}

	var skills []SkillData
	var completedQuests []string

	for _, el := range skilldata.Skillvalues {
		d := SkillData{}
		d.Level = el.Level
		d.Rank = el.Rank
		d.XP = el.Xp
		name, found := skillsIDMap[el.ID]
		if !found {
			log.Println(err)
		} else {
			d.Name = name
			skills = append(skills, d)
		}
	}
	for _, qst := range questsdata.Quests {
		if qst.Status == "COMPLETED" {
			completedQuests = append(completedQuests, qst.Title)
		}
	}

	doc, err := html.Parse(strings.NewReader(hardcoreIronmanWebpage))
	if err != nil {
		return nil, nil, false, err
	}
	//if this node is present then the search in hardcore ironman hiscores failed, the player is not hardcore ironman
	notFoundNode := FindNodeByClass(doc, "div", "tempHSUserSearchError")

	return skills, completedQuests, (notFoundNode != nil), nil
}
