package bot

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
	"gonum.org/v1/gonum/stat/combin"
	"math"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// SetDir for setting the directory for the bot
func (bot Data) SetDir() Data {
	bot.Dir, bot.Err = filepath.Abs(".")
	if bot.Err != nil {
		log.Err(bot.Err).Msg("failed to set magic-8ball directory")
	}
	return bot
}

// Start the Discord bot listener
func (bot Data) Start() {
	bot.GoBot, bot.Err = discordgo.New("Bot " + bot.Token)
	if bot.Err != nil {
		log.Err(bot.Err).Msg("failed to instantiate magic-8ball bot")
		return
	}
	bot.User, bot.Err = bot.GoBot.User(UserID)
	if bot.Err != nil {
		log.Err(bot.Err).Msg("failed to set magic-8ball user id")
		return
	}
	bot.GoBot.AddHandler(bot.MessageHandler)
	bot.GoBot.AddHandler(bot.ReactionHandler)
	bot.Err = bot.GoBot.Open()
	if bot.Err != nil {
		log.Err(bot.Err).Msg("failed to start magic-8ball listener")
		return
	}
	log.Info().Msg("magic-8ball listening")
	return
}

// MessageHandler for interpreting which function to launch from message contents
func (bot Data) MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == bot.User.ID {
		return
	}
	if strings.Contains(strings.ToLower(m.Content), "!game") {
		bot.HandleGameDay(s, m)
	}
	if strings.Contains(strings.ToLower(m.Content), "!line") {
		bot.HandleLineups(s, m)
	}
	if strings.Contains(strings.ToLower(m.Content), "!sl") {
		bot.HandleSLMatchups(s, m)
	}
	if strings.Contains(strings.ToLower(m.Content), "!inn") {
		bot.HandleHandicapAvg(s, m)
	}
	if strings.Contains(strings.ToLower(m.Content), "!opt") {
		bot.HandleOptimal(s, m)
	}
	if strings.Contains(strings.ToLower(m.Content), "!play") {
		bot.HandlePlayoff(s, m)
	}
	// inactive functions
	//if strings.Contains(strings.ToLower(m.Content), "bca") {
	//	bot.HandleBCA(s, m)
	//}
	//if containsNineBall(strings.ToLower(m.Content)) {
	//	bot.Handle9Ball(s, m)
	//}
	return
}

// ReactionHandler for interpreting how to respond to reactions
func (bot Data) ReactionHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.Member.User.ID == bot.User.ID {
		return
	}
	if slices.Contains(GameDayReactions, r.MessageReaction.Emoji.Name) &&
		(r.MessageReaction.ChannelID == GameNightChannelID ||
			r.MessageReaction.ChannelID == DevChannelID) {
		bot.HandleGameDayReaction(s, r)
	}
	return
}

// HandleGameDayReaction for handling the reaction to the game day post
func (bot Data) HandleGameDayReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	log.Info().Msg("handling reaction to game day post")
	date := time.Now()
	loc, err := time.LoadLocation("America/New_York")
	if r.ChannelID == GameNightChannelID {
		if err != nil {
			bot.Err = err
			log.Err(bot.Err).Msg("failed to load timezone, using UTC as EST+5")
			if date.Weekday() != time.Tuesday {
				log.Info().Msg("not Tuesday UTC, ignoring reaction")
				return
			}
		} else {
			date = date.In(loc)
			if date.Weekday() != time.Tuesday || date.Hour() >= 19 {
				log.Info().Msg("not Tuesday before 7pm EST, ignoring reaction")
				return
			}
		}
	}
	var status string
	switch r.MessageReaction.Emoji.Name {
	case "👍":
		status = "available"
	case "👎":
		status = "unavailable"
	case "⌛":
		status = "late"
	case "⏳":
		status = "late"
	default:
		log.Info().Msg("unknown reaction")
		return
	}
	name := r.Member.Nick
	if name == "" {
		name = r.Member.User.Username
	}
	message := discordgo.MessageSend{
		Content: fmt.Sprintf(
			"%s will be %s tonight.", name, status,
		),
	}
	if r.MessageReaction.ChannelID == DevChannelID {
		_, bot.Err = s.ChannelMessageSendComplex(DevChannelID, &message)
	} else {
		_, bot.Err = s.ChannelMessageSendComplex(GameNightChannelID, &message)
	}
	if bot.Err != nil {
		log.Err(bot.Err).Msg("failed to post message")
		return
	}
	log.Info().Msgf("%s reaction from %s to game day announcement posted to Discord channel %s",
		r.MessageReaction.Emoji.Name, name, r.ChannelID)
	return
}

// HandleGameDay for posting game day message
func (bot Data) HandleGameDay(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info().Msg("handling game day post creation")
	var opponentTeam string
	for i, name := range DivisionTeamNames {
		for _, junk := range []string{"'", "-", "8"} {
			name = strings.Replace(name, junk, "", -1)
		}
		if strings.Contains(strings.ToLower(m.Content), strings.ToLower(name)) {
			opponentTeam = DivisionTeamNames[i]
		}
	}
	message := discordgo.MessageSend{
		Content: fmt.Sprintf(
			"@everyone It's Game Day! Tonight we play %s.\n"+
				ReactionRequest, opponentTeam,
		),
	}
	if m.ChannelID == DevChannelID {
		_, bot.Err = s.ChannelMessageSendComplex(DevChannelID, &message)
	} else {
		_, bot.Err = s.ChannelMessageSendComplex(GameNightChannelID, &message)
	}
	if bot.Err != nil {
		log.Err(bot.Err).Msg("failed to post message")
		return
	}
	log.Info().Msgf("game day vs %s posted to Discord channel %s", opponentTeam, m.ChannelID)
	return
}

// HandleLineups for returning eligible lineups from a provided list of players
func (bot Data) HandleLineups(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info().Msg("handling lineups")
	re := regexp.MustCompile("[2-7]")
	var content string
	contentSlice := strings.Split(m.Content, "!")
	for _, com := range contentSlice {
		if strings.Contains(com, "line") {
			content = com
		}
	}
	skillLevelsString := re.FindAllString(content, 8)
	skillLevels := make([]int, len(skillLevelsString))
	for i, s := range skillLevelsString {
		skillLevels[i], _ = strconv.Atoi(s)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(skillLevels)))
	message := discordgo.MessageSend{}
	var lineups [][]int
	log.Info().Msgf("generating possible lineups of %v", skillLevels)
	if len(skillLevels) >= 5 {
		gen := combin.NewPermutationGenerator(len(skillLevels), 5)
		var i int
		for gen.Next() {
			permutationIndices := gen.Permutation(nil)
			var lineup []int
			for _, permutationIndex := range permutationIndices {
				lineup = append(lineup, skillLevels[permutationIndex])
			}
			sort.Sort(sort.Reverse(sort.IntSlice(lineup)))
			if validLineup(lineup) && !containsSlice(lineups, lineup) {
				lineups = append(lineups, lineup)
			}
			i++
		}
	} else {
		lineups = append(lineups, skillLevels)
	}
	var teamLineups []TeamLineup
	for _, lineup := range lineups {
		teamLineups = append(teamLineups, TeamLineup{Lineup: lineup, Sum: sum(lineup)})
	}
	sort.Slice(teamLineups[:], func(i, j int) bool {
		return teamLineups[i].Sum > teamLineups[j].Sum
	})
	for _, teamLineup := range teamLineups {
		message.Content += fmt.Sprintf("%v %v\n", teamLineup.Lineup, teamLineup.Sum)
	}
	if len(teamLineups) == 0 {
		message.Content = "No eligible lineups found"
	}
	message.Content = "Eligible Lineups:\n```" + message.Content + "```"
	if m.ChannelID == DevChannelID {
		_, bot.Err = s.ChannelMessageSendComplex(DevChannelID, &message)
	} else {
		_, bot.Err = s.ChannelMessageSendComplex(StrategyChannelID, &message)
	}
	if bot.Err != nil {
		log.Err(bot.Err).Msg("failed to post message")
		return
	}
	log.Info().Msgf("%v possible lineups posted to Discord channel %s", len(teamLineups), m.ChannelID)
	return
}

// HandleSLMatchups for returning chart of the best skill level match-ups
func (bot Data) HandleSLMatchups(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info().Msg("handling skill level match-ups")
	bot.Excel, bot.Err = excelize.OpenFile(bot.Dir + SLMatchupFile)
	if bot.Err != nil {
		log.Err(bot.Err).Msgf("failed to read excel file \"%s\"", bot.Dir+SLMatchupFile)
		return
	}
	message := discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				URL:   SLHeatMatchupAveragesUrl,
				Type:  discordgo.EmbedTypeLink,
				Title: "Skill Level Match-Ups Averages Heatmap",
			},
			{
				URL:   SLMatchupMediansUrl,
				Type:  discordgo.EmbedTypeLink,
				Title: "Skill Level Match-Ups Medians",
			},
			{
				URL:   SLMatchupModesUrl,
				Type:  discordgo.EmbedTypeLink,
				Title: "Skill Level Match-Ups Modes",
			},
		},
		Content: "```",
	}
	bot.ExcelRows = bot.Excel.GetRows(Sheet1)
	for irow, row := range bot.ExcelRows {
		for icol, colCell := range row {
			colCell = strings.Replace(colCell, "X", "", 1)
			if icol == 0 {
				colCell += strings.Repeat(" ", 2-len(colCell))
			} else if irow == 0 {
				colCell += strings.Repeat(" ", 4-len(colCell))
			} else if len(colCell) == 3 {
				colCell += "0"
			} else if len(colCell) == 1 {
				colCell += ".00"
			}
			message.Content += colCell + " "
		}
		message.Content += "\n"
	}
	message.Content += "```"
	if m.ChannelID == DevChannelID {
		_, bot.Err = s.ChannelMessageSendComplex(DevChannelID, &message)
	} else {
		_, bot.Err = s.ChannelMessageSendComplex(StrategyChannelID, &message)
	}
	if bot.Err != nil {
		log.Err(bot.Err).Msg("failed to post message")
		return
	}
	log.Info().Msgf("skill level match-ups posted to Discord channel %s", m.ChannelID)
	return
}

// HandleHandicapAvg for returning your effective innings per game
func (bot Data) HandleHandicapAvg(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info().Msg("handling skill level match-ups")
	bot.Excel, bot.Err = excelize.OpenFile(bot.Dir + InningsFile)
	if bot.Err != nil {
		log.Err(bot.Err).Msgf("failed to read excel file \"%s\"", bot.Dir+SLMatchupFile)
		return
	}
	message := discordgo.MessageSend{Content: "```"}
	bot.ExcelRows = bot.Excel.GetRows(Sheet1)
	var longestName int
	for irow, row := range bot.ExcelRows {
		if irow > 0 {
			if len(row[0]) > longestName {
				longestName = len(row[0])
			}
			var innings []float64
			for icol, colCell := range row {
				if icol > 0 && colCell != "" {
					inning, err := strconv.ParseFloat(colCell, 64)
					if err != nil {
						bot.Err = err
						log.Err(bot.Err).Msg("failed to parse float")
					}
					innings = append(innings, inning)
				}
			}
			if len(innings) > 20 {
				innings = innings[len(innings)-20:]
			}
			sort.Float64s(innings)
			if len(innings) >= 20 {
				innings = innings[:10]
			} else if len(innings) > 2 {
				var numInnings int
				if len(innings)%2 == 0 {
					numInnings = int(math.Floor(float64(len(innings) / 2)))
				} else {
					numInnings = int(math.Floor(float64(len(innings)/2)) + 1)
				}
				innings = innings[:numInnings]
			}
			var total float64
			for _, inning := range innings {
				if inning > 10 {
					inning = 10
				}
				total += inning
			}
			average := total / float64(len(innings))
			indentSpaces := strings.Repeat(" ", longestName-len(row[0])+1)
			message.Content += fmt.Sprintf("%v%s%.2f\n", row[0], indentSpaces, average)
		}
	}
	message.Content += "```"
	if m.ChannelID == DevChannelID {
		_, bot.Err = s.ChannelMessageSendComplex(DevChannelID, &message)
	} else {
		_, bot.Err = s.ChannelMessageSendComplex(StrategyChannelID, &message)
	}
	if bot.Err != nil {
		log.Err(bot.Err).Msg("failed to post message")
		return
	}
	log.Info().Msgf("skill level match-ups posted to Discord channel %s", m.ChannelID)
	return
}

// HandleOptimal for returning max expected points lineup from opponent's lineup
func (bot Data) HandleOptimal(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info().Msg("handling optimal lineups")
	re := regexp.MustCompile("[2-7]")
	var content string
	contentSlice := strings.Split(m.Content, "!")
	for _, com := range contentSlice {
		if strings.Contains(com, "opt") {
			content = com
		}
	}
	comSlice := strings.Split(content, " ")
	if len(comSlice) < 3 {
		bot.Err = errors.New("invalid command")
		log.Err(bot.Err).Msg("not enough arguments")
		return
	}
	opponentArrString := comSlice[1]
	opponentSkillLevelsString := re.FindAllString(opponentArrString, 5)
	opponentSkillLevels := make([]int, len(opponentSkillLevelsString))
	for i, s := range opponentSkillLevelsString {
		opponentSkillLevels[i], _ = strconv.Atoi(s)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(opponentSkillLevels)))
	teamArrString := comSlice[2]
	teamSkillLevelsString := re.FindAllString(teamArrString, 8)
	teamSkillLevels := make([]int, len(teamSkillLevelsString))
	for i, s := range teamSkillLevelsString {
		teamSkillLevels[i], _ = strconv.Atoi(s)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(teamSkillLevels)))
	if len(teamSkillLevels) == 0 || len(opponentSkillLevels) != 5 {
		bot.Err = errors.New("invalid command")
		log.Err(bot.Err).Msg("not enough arguments or invalid lineup")
		return
	}
	var message discordgo.MessageSend
	var lineups [][]int
	log.Info().Msgf("generating possible lineups and expected points of %v vs %v", teamSkillLevels, opponentSkillLevels)
	if len(teamSkillLevels) >= 5 {
		gen := combin.NewPermutationGenerator(len(teamSkillLevels), 5)
		var i int
		for gen.Next() {
			permutationIndices := gen.Permutation(nil)
			var lineup []int
			for _, permutationIndex := range permutationIndices {
				lineup = append(lineup, teamSkillLevels[permutationIndex])
			}
			sort.Sort(sort.Reverse(sort.IntSlice(lineup)))
			if validLineup(lineup) && !containsSlice(lineups, lineup) {
				lineups = append(lineups, lineup)
			}
			i++
		}
	} else {
		lineups = append(lineups, teamSkillLevels)
	}
	bot.Excel, bot.Err = excelize.OpenFile(bot.Dir + SLMatchupFile)
	if bot.Err != nil {
		log.Err(bot.Err).Msgf("failed to read excel file \"%s\"", bot.Dir+SLMatchupFile)
		return
	}
	bot.ExcelRows = bot.Excel.GetRows(Sheet1)
	var teamLineups []TeamLineup
	for _, lineup := range lineups {
		var matchups []Matchup
		for _, opponentPlayer := range opponentSkillLevels {
			for _, teamPlayer := range lineup {
				points, err := strconv.ParseFloat(
					strings.Replace(bot.ExcelRows[teamPlayer-1][opponentPlayer-1],
						"X", "", 1), 64)
				if err != nil {
					bot.Err = err
					log.Err(bot.Err).Msg("failed to parse float")
					return
				}
				matchup := Matchup{
					SkillLevels:       [2]int{teamPlayer, opponentPlayer},
					ExpectedPointsFor: points,
				}
				matchups = append(matchups, matchup)
			}
		}
		teamLineups = append(teamLineups, TeamLineup{
			Lineup:         lineup,
			OpponentLineup: opponentSkillLevels,
			Matchups:       matchups,
		})
	}
	var t []TeamLineup
	for _, tl := range teamLineups {
		gen := combin.NewCombinationGenerator(len(tl.Matchups), 5)
		var i int
		for gen.Next() {
			combinationIndices := gen.Combination(nil)
			var tp []int
			var op []int
			var total float64
			for _, combinationIndex := range combinationIndices {
				tp = append(tp, tl.Matchups[combinationIndex].SkillLevels[0])
				op = append(op, tl.Matchups[combinationIndex].SkillLevels[1])
				total += tl.Matchups[combinationIndex].ExpectedPointsFor
			}
			sort.Sort(sort.Reverse(sort.IntSlice(tp)))
			sort.Sort(sort.Reverse(sort.IntSlice(op)))
			if reflect.DeepEqual(tp, tl.Lineup) &&
				reflect.DeepEqual(op, tl.OpponentLineup) {
				var tls TeamLineup
				tls.MatchupExpectedPointsFor = total
				for _, combinationIndex := range combinationIndices {
					m := Matchup{
						SkillLevels:       tl.Matchups[combinationIndex].SkillLevels,
						ExpectedPointsFor: tl.Matchups[combinationIndex].ExpectedPointsFor,
					}
					tls.Matchups = append(tls.Matchups, m)
				}
				toadd := true
				for _, v := range t {
					if reflect.DeepEqual(v, tls) {
						toadd = false
					}
				}
				if toadd {
					t = append(t, tls)
				}
			}
			i++
		}
	}
	sort.Slice(t, func(i, j int) bool {
		return t[i].MatchupExpectedPointsFor > t[j].MatchupExpectedPointsFor
	})
	if len(t) > 50 {
		t = t[:50]
	}
	for _, tl := range t {
		sort.Slice(tl.Matchups, func(i, j int) bool {
			return tl.Matchups[i].SkillLevels[0] > tl.Matchups[j].SkillLevels[0]
		})
	}
	var tli []TeamLineup
	var prevm []Matchup
	for _, v := range t {
		if !reflect.DeepEqual(v.Matchups, prevm) {
			tli = append(tli, v)
			prevm = v.Matchups
		}
	}
	if len(tli) > 10 {
		tli = tli[:10]
	}
	for _, l := range tli {
		var sls string
		for _, m := range l.Matchups {
			sls += fmt.Sprintf("%d/%d ", m.SkillLevels[0], m.SkillLevels[1])
		}
		if !strings.Contains(message.Content, sls) {
			message.Content += fmt.Sprintf("%s\t%.2f\n", sls, l.MatchupExpectedPointsFor)
		}
	}
	if len(teamLineups) == 0 {
		message.Content = "No eligible lineups found"
	}
	message.Content = "Expected Points by Matchups:\n```" + message.Content + "```"
	if m.ChannelID == DevChannelID {
		_, bot.Err = s.ChannelMessageSendComplex(DevChannelID, &message)
	} else {
		_, bot.Err = s.ChannelMessageSendComplex(StrategyChannelID, &message)
	}
	if bot.Err != nil {
		log.Err(bot.Err).Msg("failed to post message")
		return
	}
	log.Info().Msgf("top possible lineup and expected points posted to Discord channel %s", m.ChannelID)
	return
}

// HandlePlayoff for returning max differential expected points lineup from opponent's lineup
func (bot Data) HandlePlayoff(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info().Msg("handling optimal playoff lineups")
	re := regexp.MustCompile("[2-7]")
	var content string
	contentSlice := strings.Split(m.Content, "!")
	for _, com := range contentSlice {
		if strings.Contains(com, "play") {
			content = com
		}
	}
	comSlice := strings.Split(content, " ")
	if len(comSlice) < 3 {
		bot.Err = errors.New("invalid command")
		log.Err(bot.Err).Msg("not enough arguments")
		return
	}
	opponentArrString := comSlice[1]
	opponentSkillLevelsString := re.FindAllString(opponentArrString, 5)
	opponentSkillLevels := make([]int, len(opponentSkillLevelsString))
	for i, s := range opponentSkillLevelsString {
		opponentSkillLevels[i], _ = strconv.Atoi(s)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(opponentSkillLevels)))
	teamArrString := comSlice[2]
	teamSkillLevelsString := re.FindAllString(teamArrString, 8)
	teamSkillLevels := make([]int, len(teamSkillLevelsString))
	for i, s := range teamSkillLevelsString {
		teamSkillLevels[i], _ = strconv.Atoi(s)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(teamSkillLevels)))
	if len(teamSkillLevels) == 0 || len(opponentSkillLevels) != 5 {
		bot.Err = errors.New("invalid command")
		log.Err(bot.Err).Msg("not enough arguments or invalid lineup")
		return
	}
	var message discordgo.MessageSend
	var lineups [][]int
	log.Info().Msgf("generating possible lineups and expected points of %v vs %v", teamSkillLevels, opponentSkillLevels)
	if len(teamSkillLevels) >= 5 {
		gen := combin.NewPermutationGenerator(len(teamSkillLevels), 5)
		var i int
		for gen.Next() {
			permutationIndices := gen.Permutation(nil)
			var lineup []int
			for _, permutationIndex := range permutationIndices {
				lineup = append(lineup, teamSkillLevels[permutationIndex])
			}
			sort.Sort(sort.Reverse(sort.IntSlice(lineup)))
			if validLineup(lineup) && !containsSlice(lineups, lineup) {
				lineups = append(lineups, lineup)
			}
			i++
		}
	} else {
		lineups = append(lineups, teamSkillLevels)
	}
	bot.Excel, bot.Err = excelize.OpenFile(bot.Dir + SLMatchupFile)
	if bot.Err != nil {
		log.Err(bot.Err).Msgf("failed to read excel file \"%s\"", bot.Dir+SLMatchupFile)
		return
	}
	bot.ExcelRows = bot.Excel.GetRows(Sheet1)
	var teamLineups []TeamLineup
	for _, lineup := range lineups {
		var matchups []Matchup
		for _, opponentPlayer := range opponentSkillLevels {
			for _, teamPlayer := range lineup {
				teamPoints, err := strconv.ParseFloat(
					strings.Replace(bot.ExcelRows[teamPlayer-1][opponentPlayer-1],
						"X", "", 1), 64)
				if err != nil {
					bot.Err = err
					log.Err(bot.Err).Msg("failed to parse float")
					return
				}
				opponentPoints, err := strconv.ParseFloat(
					strings.Replace(bot.ExcelRows[opponentPlayer-1][teamPlayer-1],
						"X", "", 1), 64)
				if err != nil {
					bot.Err = err
					log.Err(bot.Err).Msg("failed to parse float")
					return
				}
				matchup := Matchup{
					SkillLevels:           [2]int{teamPlayer, opponentPlayer},
					ExpectedPointsFor:     teamPoints,
					ExpectedPointsAgainst: opponentPoints,
				}
				matchups = append(matchups, matchup)
			}
		}
		teamLineups = append(teamLineups, TeamLineup{
			Lineup:         lineup,
			OpponentLineup: opponentSkillLevels,
			Matchups:       matchups,
		})
	}
	var t []TeamLineup
	for _, tl := range teamLineups {
		gen := combin.NewCombinationGenerator(len(tl.Matchups), 5)
		var i int
		for gen.Next() {
			combinationIndices := gen.Combination(nil)
			var tp []int
			var op []int
			var teamTotal float64
			var opponentTotal float64
			var diffTotal float64
			for _, combinationIndex := range combinationIndices {
				tp = append(tp, tl.Matchups[combinationIndex].SkillLevels[0])
				op = append(op, tl.Matchups[combinationIndex].SkillLevels[1])
				teamTotal += tl.Matchups[combinationIndex].ExpectedPointsFor
				opponentTotal += tl.Matchups[combinationIndex].ExpectedPointsAgainst
				diffTotal += tl.Matchups[combinationIndex].ExpectedPointsFor - tl.Matchups[combinationIndex].ExpectedPointsAgainst
			}
			sort.Sort(sort.Reverse(sort.IntSlice(tp)))
			sort.Sort(sort.Reverse(sort.IntSlice(op)))
			if reflect.DeepEqual(tp, tl.Lineup) &&
				reflect.DeepEqual(op, tl.OpponentLineup) {
				var tls TeamLineup
				tls.MatchupExpectedPointsFor = teamTotal
				tls.MatchupExpectedPointsAgainst = opponentTotal
				tls.MatchupExpectedPointsDifference = diffTotal
				for _, combinationIndex := range combinationIndices {
					m := Matchup{
						SkillLevels:              tl.Matchups[combinationIndex].SkillLevels,
						ExpectedPointsFor:        tl.Matchups[combinationIndex].ExpectedPointsFor,
						ExpectedPointsAgainst:    tl.Matchups[combinationIndex].ExpectedPointsAgainst,
						ExpectedPointsDifference: tl.Matchups[combinationIndex].ExpectedPointsFor - tl.Matchups[combinationIndex].ExpectedPointsAgainst,
					}
					tls.Matchups = append(tls.Matchups, m)
				}
				toadd := true
				for _, v := range t {
					if reflect.DeepEqual(v, tls) {
						toadd = false
					}
				}
				if toadd {
					t = append(t, tls)
				}
			}
			i++
		}
	}
	sort.Slice(t, func(i, j int) bool {
		return t[i].MatchupExpectedPointsDifference > t[j].MatchupExpectedPointsDifference
	})
	if len(t) > 50 {
		t = t[:50]
	}
	for _, tl := range t {
		sort.Slice(tl.Matchups, func(i, j int) bool {
			return tl.Matchups[i].SkillLevels[0] > tl.Matchups[j].SkillLevels[0]
		})
	}
	var tli []TeamLineup
	var prevm []Matchup
	for _, v := range t {
		if !reflect.DeepEqual(v.Matchups, prevm) {
			tli = append(tli, v)
			prevm = v.Matchups
		}
	}
	if len(tli) > 10 {
		tli = tli[:10]
	}
	for _, l := range tli {
		var sls string
		for i, m := range l.Matchups {
			sls += fmt.Sprintf("%d/%d", m.SkillLevels[0], m.SkillLevels[1])
			if i != len(l.Matchups)-1 {
				sls += "|"
			}
		}
		if !strings.Contains(message.Content, sls) {
			message.Content += fmt.Sprintf("%s|%.2f/%.2f|%.2f\n", sls, l.MatchupExpectedPointsFor, l.MatchupExpectedPointsAgainst, l.MatchupExpectedPointsDifference)
		}
	}
	if len(teamLineups) == 0 {
		message.Content = "No eligible lineups found"
	}
	message.Content = "Expected Points Favored by Matchups:\n```M1  M2  M3  M4  M5  ExpPoints Diff\n" + message.Content + "```"
	if m.ChannelID == DevChannelID {
		_, bot.Err = s.ChannelMessageSendComplex(DevChannelID, &message)
	} else {
		_, bot.Err = s.ChannelMessageSendComplex(StrategyChannelID, &message)
	}
	if bot.Err != nil {
		log.Err(bot.Err).Msg("failed to post message")
		return
	}
	log.Info().Msgf("top possible lineup and expected points posted to Discord channel %s", m.ChannelID)
	return
}

// HandleBCA for mentions of BCA play
func (bot Data) HandleBCA(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info().Msg("handling mention of BCA play")
	message := discordgo.MessageSend{
		Content: "BCA is for bums.",
		TTS:     true,
	}
	_, bot.Err = s.ChannelMessageSendComplex(m.ChannelID, &message)
	if bot.Err != nil {
		log.Err(bot.Err).Msg("failed to post message")
		return
	}
	log.Info().Msgf("bca rebuttal posted to %s in Discord channel %s", m.Member.Nick, m.ChannelID)
	return
}

// Handle9Ball for mentions of 9-ball play
func (bot Data) Handle9Ball(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info().Msg("handling mention of 9-ball play")
	message := discordgo.MessageSend{
		Content: "9-Ball is for bums.",
		TTS:     true,
	}
	_, bot.Err = s.ChannelMessageSendComplex(m.ChannelID, &message)
	if bot.Err != nil {
		log.Err(bot.Err).Msg("failed to post message")
		return
	}
	log.Info().Msgf("9-ball rebuttal posted to %s in Discord channel %s", m.Member.Nick, m.ChannelID)
	return
}

// seniorSkillRule returns a bool indicating if a lineup violates this rule
func seniorSkillRule(array []int) bool {
	numSeniors := 0
	for _, v := range array {
		if v >= SeniorSkillLevel {
			numSeniors++
		}
	}
	if numSeniors > 2 {
		return true
	}
	return false
}

// sum returns the sum of the elements in the given int slice
func sum(array []int) int {
	result := 0
	for _, v := range array {
		result += v
	}
	return result
}

// containsSlice returns true if the given slice of int slices contains the given int slice
func containsSlice(slc [][]int, ele []int) bool {
	for _, v := range slc {
		if reflect.DeepEqual(v, ele) {
			return true
		}
	}
	return false
}

// validLineup returns true if the given lineup is valid
func validLineup(lineup []int) bool {
	if len(lineup) != 5 {
		return false
	}
	if sum(lineup) < 10 || sum(lineup) > 23 {
		return false
	}
	if seniorSkillRule(lineup) {
		return false
	}
	return true
}

// containsNineBall returns true if the message mentions 9-ball
func containsNineBall(text string) bool {
	return strings.Contains(strings.Replace(strings.Replace(
		strings.ToLower(text), "-", "", -1), " ", "", -1), "9ball") ||
		strings.Contains(strings.Replace(strings.ToLower(text), " ", "", -1), "nineball")
}
