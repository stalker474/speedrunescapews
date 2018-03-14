package main

//SupportedGames List of currently supported games
var SupportedGames = [...]JSONSupportedGame{
	JSONSupportedGame{
		"Runescape3",
		true,
		true,
		true,
		true,
		true,
		false,
		true,
		"www.runescape.com",
		"secure.runescape.com/m=account-creation/l=0/create_account",
		"https://www.runescape.com/webfiles/latest/common/img/logos/runescape.png",
		[]JSONSupportedGameChallenge{
			{"Total Lvl", "http://vignette.wikia.nocookie.net/ikov-2/images/2/25/Unnamed_%281%29.png"},
			{"Combat Lvl", "http://vignette.wikia.nocookie.net/runescape2/images/0/00/Attack.png"},
			{"Complete Demon Slayer quest", "http://vignette.wikia.nocookie.net/runescape2/images/8/8d/Quest_Icon_Crest.png"},
			{"Complete Dragon Slayer quest", "http://vignette.wikia.nocookie.net/runescape2/images/8/8d/Quest_Icon_Crest.png"},
		},
	},
	JSONSupportedGame{
		"Old School Runescape",
		true,
		false,
		true,
		true,
		false,
		false,
		false,
		"oldschool.runescape.com/",
		"secure.runescape.com/m=account-creation/g=oldscape/create_account",
		"https://oldschool.runescape.com/webfiles/latest/common/img/logos/oldschool.png",
		[]JSONSupportedGameChallenge{
			{"Total Lvl", "http://vignette.wikia.nocookie.net/ikov-2/images/2/25/Unnamed_%281%29.png"},
			{"Combat Lvl", "http://vignette.wikia.nocookie.net/runescape2/images/0/00/Attack.png"},
		},
	},
}
