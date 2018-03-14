package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//Useful urls

var wsURLGetUltimateIronman = "http://services.runescape.com/m=hiscore_oldschool_ultimate/hiscorepersonal.ws?user1="
var wsURLGetHardcoreIronman = "http://services.runescape.com/m=hiscore_oldschool_hardcore_ironman/hiscorepersonal.ws?user1="
var wsURLGetNormal = "http://services.runescape.com/m=hiscore_oldschool/hiscorepersonal.ws?user1="

var wsURLGetStatsRS3 = "https://apps.runescape.com/runemetrics/app/levels/player/"
var wssURLGetStatsApp = "https://apps.runescape.com/runemetrics/profile/profile?activities=20&user="
var wssURLGetQuestsApp = "https://apps.runescape.com/runemetrics/quests?user="
var wsURLGetUserHardcoreIronmanRS3 = "http://services.runescape.com/m=hiscore_hardcore_ironman/ranking?user="

// HiscoresData blablabla
type HiscoresData struct {
	Data            map[string]SkillData
	CompletedQuests []string
	Status          IronmanStatus
}

// SkillData data of a single skill
type SkillData struct {
	Name  string
	Rank  int
	Level int
	XP    int
}

// IronmanStatus blablabla
type IronmanStatus int

// Ironman types
const (
	NORMAL   = IronmanStatus(0)
	HARDCORE = IronmanStatus(1)
	ULTIMATE = IronmanStatus(2)
)

// Hiscores Data fetcher for RS hiscores
type Hiscores struct {
}

// Retrieve get a player hiscore by name and ironman status
func (*Hiscores) Retrieve(PlayerName string, Status IronmanStatus, rs3 bool) (*HiscoresData, error) {
	hc := http.Client{}
	var url string

	if !rs3 {
		switch st := Status; st {
		case NORMAL:
			url = wsURLGetNormal
		case HARDCORE:
			url = wsURLGetHardcoreIronman
		case ULTIMATE:
			url = wsURLGetUltimateIronman
		default:
			return nil, errors.New("unknown status")
		}
	} else {
		url = wssURLGetStatsApp
	}

	playerNameUrled := strings.Replace(PlayerName, " ", "%A0", -1)

	req, _ := http.NewRequest("GET", url+playerNameUrled, strings.NewReader(""))

	resp, err := hc.Do(req)

	if resp.StatusCode != 200 {
		return nil, errors.New("Unable to reach stats page")
	}

	if err != nil {
		log.Print(err)
		return nil, err
	}
	var scores []SkillData
	var quests []string
	var isHardcore bool
	if !rs3 {
		osrs := OSRS{}
		scores, err = osrs.Retrieve(resp.Body)
		isHardcore = true
	} else {
		rs3 := RS3{}
		body, _ := ioutil.ReadAll(resp.Body)

		//now get quests data
		req, _ = http.NewRequest("GET", wssURLGetQuestsApp+playerNameUrled, strings.NewReader(""))
		resp, err = hc.Do(req)

		if resp.StatusCode != 200 {
			return nil, errors.New("Unable to reach quests page")
		}

		bodyQuests, _ := ioutil.ReadAll(resp.Body)

		//now get hardcore ironman data
		req, _ = http.NewRequest("GET", wsURLGetUserHardcoreIronmanRS3+playerNameUrled, strings.NewReader(""))
		resp, err = hc.Do(req)

		if resp.StatusCode != 200 {
			return nil, errors.New("Unable to reach hardcore ironman hiscores page")
		}

		hardcoreIronmanHiscoresWebpage, _ := ioutil.ReadAll(resp.Body)

		scores, quests, isHardcore, err = rs3.Retrieve(string(body), string(bodyQuests), string(hardcoreIronmanHiscoresWebpage))
	}

	if err != nil {
		log.Print(err)
		return nil, err
	}

	//log.Println(str)
	res := new(HiscoresData)
	res.Data = make(map[string]SkillData)
	res.CompletedQuests = quests
	if isHardcore {
		res.Status = HARDCORE
	} else {
		res.Status = NORMAL
	}

	for _, el := range scores {
		res.Data[el.Name] = el
	}
	return res, nil
}
