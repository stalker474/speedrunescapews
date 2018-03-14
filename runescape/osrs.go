package runescape

import (
	"errors"
	"io"
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

//OSRS data fetcher for runescape
type OSRS struct {
}

// getTable find the contentHiscore table in the dom
func getTable(doc *html.Node) (*html.Node, error) {
	hiscoresContentDiv := FindNode(doc, "div", "contentHiscores")
	if hiscoresContentDiv != nil {
		table := FindNode(hiscoresContentDiv, "tbody", "")
		if table != nil {
			return table, nil
		}
		return nil, errors.New("Missing tbody in the hiscores div")
	}

	return nil, errors.New("Missing <div id=\"contentHiscores\"> in the node tree")
}

// getScoresList retrieve scores as a map of string,SkillData
func getScoresOSRSList(n *html.Node) ([]SkillData, error) {
	var trList []*html.Node
	var skills []SkillData
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		found := FindNode(c, "a", "")
		if found != nil {
			trList = append(trList, c)
		}
	}
	/*
		Each element of this list will look like that
		   <tr>
		   <td align="right"><img class="miniimg" src="http://www.runescape.com/img/rsp777/hiscores/skill_icon_attack1.gif"></td>
		   <td align="left"><a href="overall.ws?table=1&amp;user=Test474">
		   Attack
		   </a></td>
		   <td align="right">38,496</td>
		   <td align="right">1</td>
		   <td align="right">12</td>

		   </tr>
	*/
	for _, element := range trList {
		var skill SkillData
		a := FindNode(element, "a", "")

		tdRank := a.Parent.NextSibling.NextSibling
		if tdRank == nil {
			log.Println("Unable to find rank td")
			log.Print(RenderNode(element))
			continue
		}
		tdLevel := tdRank.NextSibling.NextSibling
		if tdLevel == nil {
			log.Println("Unable to find level td")
			log.Print(RenderNode(element))
			continue
		}
		tdXP := tdLevel.NextSibling.NextSibling
		if tdXP == nil {
			log.Println("Unable to find xp td")
			log.Print(RenderNode(element))
			continue
		}

		rankContent := strings.Replace(tdRank.FirstChild.Data, ",", "", -1)
		levelContent := strings.Replace(tdLevel.FirstChild.Data, ",", "", -1)
		xpContent := strings.Replace(tdXP.FirstChild.Data, ",", "", -1)

		rank, err := strconv.Atoi(rankContent)
		if err != nil {
			log.Println("Failed to parse rank value in td")
			log.Print(RenderNode(element))
			continue
		}
		level, err := strconv.Atoi(levelContent)
		if err != nil {
			log.Println("Failed to parse level value in td")
			log.Print(RenderNode(element))
			continue
		}
		xp, err := strconv.Atoi(xpContent)
		if err != nil {
			log.Println("Failed to parse xp value in td")
			log.Print(RenderNode(element))
			continue
		}
		skill.Level = level
		skill.Rank = rank
		skill.XP = xp

		if a.FirstChild != nil {
			skill.Name = strings.Trim(a.FirstChild.Data, "\n")
			skills = append(skills, skill)
		} else {
			log.Println("Error processing skill tr")
			log.Print(RenderNode(element))
			continue
		}
	}

	if len(skills) == 0 {
		return nil, errors.New("Unable to find skills in table")
	}

	return skills, nil
}

// Retrieve Simply fetch the skill data from dom
func (*OSRS) Retrieve(r io.Reader) ([]SkillData, error) {

	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	table, err := getTable(doc)
	if err != nil {
		return nil, err
	}

	scores, err := getScoresOSRSList(table)
	if err != nil {
		return nil, err
	}
	return scores, nil
}
