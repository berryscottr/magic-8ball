package bot

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
	"gonum.org/v1/gonum/stat/combin"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// MessageHandler for interpreting which function to launch from message contents
func (bot *Data) MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == bot.User.ID {
		return
	}
	if strings.Contains(strings.ToLower(m.Content), "!8game") {
		bot.HandleGameDay(s, m, WookieMistakes.Name)
	}
	if strings.Contains(strings.ToLower(m.Content), "!9game") {
		bot.HandleGameDay(s, m, SafetyDance.Name)
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
	if strings.Contains(strings.ToLower(m.Content), "!opt8") {
		bot.HandleOptimal8(s, m)
	}
	if strings.Contains(strings.ToLower(m.Content), "!opt9") {
		bot.HandleOptimal9(s, m)
	}
	if strings.Contains(strings.ToLower(m.Content), "!play") {
		bot.HandlePlayoff(s, m)
	}
	if strings.Contains(strings.ToLower(m.Content), "!cal") {
		bot.HandleCalendar(s, m)
	}
}

// ReactionHandler for interpreting how to respond to reactions
func (bot *Data) ReactionHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.Member.User.ID == bot.User.ID {
		return
	}
	if slices.Contains(GameDayReactions, r.MessageReaction.Emoji.Name) &&
		(r.MessageReaction.ChannelID == GameNight8ChannelID || r.MessageReaction.ChannelID == GameNight9ChannelID || r.MessageReaction.ChannelID == DevChannelID) {
		bot.HandleGameDayReaction(s, r)
	}
}

// HandleGameDayReaction for handling the reaction to the game day post
func (bot *Data) HandleGameDayReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	log.Info().Msg("handling reaction to game day post")
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
	var teammate Teammate
	for _, t := range Teammates {
		if t.UserID == r.MessageReaction.UserID {
			teammate = t
			log.Info().Msgf("tracking reaction from teammate: %s", teammate.LastName)
			break
		}
	}
	if teammate.UserID == "" {
		log.Info().Msg("unknown teammate")
		return
	}
	var oldMsg *discordgo.Message
	oldMsg, bot.Err = s.ChannelMessage(r.MessageReaction.ChannelID, r.MessageReaction.MessageID)
	if bot.Err != nil {
		log.Err(bot.Err).Msgf("failed to find reaction message in Discord channel %s", r.MessageReaction.ChannelID)
	}
	var newMsg string
	var newLine string
	msgLines := strings.Split(oldMsg.Content, "\n")
	for _, line := range msgLines {
		if strings.Contains(line, teammate.LastName) {
			log.Info().Msgf("modifying attendance for teammate: %s", teammate.LastName)
			switch status {
			case "available":
				newLine = strings.Replace(line, "|     |", "|  X  |", 1)
				newLine = strings.Replace(newLine, "|  X   |", "|      |", 1)
				newLine = strings.Replace(newLine, "| X  |", "|    |", 1)
			case "unavailable":
				newLine = strings.Replace(line, "|  X  |", "|     |", 1)
				newLine = strings.Replace(newLine, "|  X   |", "|      |", 1)
				newLine = strings.Replace(newLine, "|    |", "| X  |", 1)
			case "late":
				newLine = strings.Replace(line, "|  X  |", "|     |", 1)
				newLine = strings.Replace(newLine, "|      |", "|  X   |", 1)
				newLine = strings.Replace(newLine, "| X  |", "|    |", 1)
			}
			newMsg = strings.Replace(oldMsg.Content, line, newLine, 1)
			break
		}
	}
	_, bot.Err = s.ChannelMessageEdit(r.MessageReaction.ChannelID, r.MessageReaction.MessageID, newMsg)
	if bot.Err != nil {
		log.Err(bot.Err).Msgf("failed to edit message in Discord channel %s", r.MessageReaction.ChannelID)
	} else {
		log.Info().Msgf("%s reaction from %s to game day announcement posted in Discord channel %s has been updated",
			r.MessageReaction.Emoji.Name, teammate.LastName, r.MessageReaction.ChannelID)
	}
}

// HandleGameDay for posting game day message
func (bot *Data) HandleGameDay(s *discordgo.Session, m *discordgo.MessageCreate, teamName string) {
	log.Info().Msg("handling game day post creation")
	var team Team
	if teamName == WookieMistakes.Name {
		team = WookieMistakes
	} else if teamName == SafetyDance.Name {
		team = SafetyDance
	} else {
		bot.Err = errors.New("invalid team name")
		log.Err(bot.Err).Msgf("failed to create game day post for team: %s", teamName)
		return
	}
	var opponentTeam string
	for i, name := range team.DivisionTeamNames {
		for _, junk := range []string{"'", "-", "8", "9"} {
			name = strings.Replace(name, junk, "", -1)
		}
		name = strings.TrimRight(name, " ")
		if strings.Contains(strings.ToLower(m.Content), strings.ToLower(name)) {
			opponentTeam = team.DivisionTeamNames[i]
		}
	}
	var customMessage string
	if strings.Count(m.Content, "\"") == 2 {
		customMessage = fmt.Sprintf("\n\n\"%s\" - %s", strings.Split(m.Content, "\"")[1], m.Author.Username)
	}
	message := discordgo.MessageSend{
		Content: fmt.Sprintf(
			"@everyone It's Game Day! Tonight %s plays %s.\n"+
				ReactionRequest+customMessage, team.Name, opponentTeam,
		),
	}
	message.Content += "\n```\n"
	message.Content += "+---Name---+-Yes-+-Late-+-No-+\n"
	var longestName int
	for _, teammate := range Teammates {
		for _, t := range teammate.Teams {
			if t.Name == team.Name {
				if len(teammate.LastName) > longestName {
					longestName = len(teammate.LastName)
				}
			}
		}
	}
	var numspaces int
	for _, teammate := range Teammates {
		for _, t := range teammate.Teams {
			if t.Name == team.Name {
				numspaces = longestName + 1 - len(teammate.LastName)
				message.Content += fmt.Sprintf("| %s%s|     |      |    |\n", teammate.LastName, strings.Repeat(" ", numspaces))
			}
		}
	}
	message.Content += "+----------+-----+------+----+\n```"
	if m.ChannelID == DevChannelID {
		if strings.Contains(m.Content, "--now") {
			_, bot.Err = s.ChannelMessageSendComplex(DevChannelID, &message)
		} else if strings.Contains(m.Content, "--schedule") {
			time.AfterFunc(time.Minute, func() {
				_, bot.Err = s.ChannelMessageSendComplex(DevChannelID, &message)
			})
		}	else {
			_, bot.Err = s.ChannelMessageSendComplex(DevChannelID, &message)
		}
	} else {
		loc, err := time.LoadLocation("America/New_York")
		if err != nil {
			loc = time.UTC
			log.Err(bot.Err).Msg("failed to load timezone, using UTC as EST+5")
		}
		date := time.Now().In(loc)
		if strings.Contains(m.Content, "--now") {
			_, bot.Err = s.ChannelMessageSendComplex(team.GameNightChannelID, &message)
		} else if date.Weekday() != time.Tuesday {
			nextTuesday := date.AddDate(0, 0, int((time.Tuesday - date.Weekday() + 7) % 7))
			scheduleTime := time.Date(date.Year(), date.Month(), nextTuesday.Day(), 8, 55, 0, 0, loc)
			time.AfterFunc(scheduleTime.Sub(date), func() {
				_, bot.Err = s.ChannelMessageSendComplex(team.GameNightChannelID, &message)
			})
		} else if date.Hour() < 5 {
			scheduleTime := time.Date(date.Year(), date.Month(), date.Day(), 8, 55, 0, 0, loc)
			time.AfterFunc(scheduleTime.Sub(date), func() {
				_, bot.Err = s.ChannelMessageSendComplex(team.GameNightChannelID, &message)
			})
		} else {
			_, bot.Err = s.ChannelMessageSendComplex(team.GameNightChannelID, &message)
		}
	}
	if bot.Err != nil {
		log.Err(bot.Err).Msg("failed to post message")
		return
	}
	log.Info().Msgf("game day %s vs %s posted or scheduled to Discord channel %s", team.Name, opponentTeam, m.ChannelID)
}

// HandleLineups for returning eligible lineups from a provided list of players
func (bot *Data) HandleLineups(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info().Msg("handling lineups")
	re := regexp.MustCompile("[1-9]")
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
}

// HandleSLMatchups for returning chart of the best skill level match-ups
func (bot *Data) HandleSLMatchups(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info().Msg("handling skill level match-ups")
	bot.Excel, bot.Err = excelize.OpenFile(bot.Dir + SLMatchupFile)
	if bot.Err != nil {
		log.Err(bot.Err).Msgf("failed to read excel file \"%s\"", bot.Dir+SLMatchupFile)
		return
	}
	message := discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				URL:   EightSLHeatMatchupAveragesUrl,
				Type:  discordgo.EmbedTypeLink,
				Title: "8-Ball Skill Level Match-Up Averages Heatmap",
			},
			{
				URL:   NineSLHeatMatchupAveragesUrl,
				Type:  discordgo.EmbedTypeLink,
				Title: "9-Ball Skill Level Match-Up Averages Heatmap",
			},
			{
				URL:   SLMatchupMediansUrl,
				Type:  discordgo.EmbedTypeLink,
				Title: "8-Ball Skill Level Match-Up Medians",
			},
			{
				URL:   SLMatchupModesUrl,
				Type:  discordgo.EmbedTypeLink,
				Title: "8-Ball Skill Level Match-Up Modes",
			},
		},
		Content: "```",
	}
	bot.ExcelRows = bot.Excel.GetRows(bot.Excel.WorkBook.Sheets.Sheet[0].Name)
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
}

// HandleHandicapAvg for returning your effective innings per game
func (bot *Data) HandleHandicapAvg(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info().Msg("handling skill level match-ups")
	bot.Excel, bot.Err = excelize.OpenFile(bot.Dir + InningsFile)
	if bot.Err != nil {
		log.Err(bot.Err).Msgf("failed to read excel file \"%s\"", bot.Dir+InningsFile)
		return
	}
	message := discordgo.MessageSend{Content: "```"}
	bot.ExcelRows = bot.Excel.GetRows(bot.Excel.WorkBook.Sheets.Sheet[0].Name)
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
					numInnings = int(float64(len(innings) / 2))
				} else {
					numInnings = int(float64(len(innings)/2)) + 1
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
	log.Info().Msgf("effective innings posted to Discord channel %s", m.ChannelID)
}

// HandleOptimal8 for returning max expected points lineup from opponent's lineup for eight-ball
func (bot *Data) HandleOptimal8(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info().Msg("handling optimal lineups")
	re := regexp.MustCompile("[2-7]")
	var content string
	contentSlice := strings.Split(m.Content, "!")
	for _, com := range contentSlice {
		if strings.Contains(com, "opt8") {
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
	bot.ExcelRows = bot.Excel.GetRows(bot.Excel.WorkBook.Sheets.Sheet[0].Name)
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
}

// HandleOptimal9 for returning max expected points lineup from opponent's lineup for nine-ball
func (bot *Data) HandleOptimal9(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info().Msg("handling optimal lineups")
	re := regexp.MustCompile("[1-9]")
	var content string
	contentSlice := strings.Split(m.Content, "!")
	for _, com := range contentSlice {
		if strings.Contains(com, "opt9") {
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
	bot.Excel, bot.Err = excelize.OpenFile(bot.Dir + SLMatchupFileNine)
	if bot.Err != nil {
		log.Err(bot.Err).Msgf("failed to read excel file \"%s\"", bot.Dir+SLMatchupFileNine)
		return
	}
	bot.ExcelRows = bot.Excel.GetRows(bot.Excel.WorkBook.Sheets.Sheet[0].Name)
	var teamLineups []TeamLineup
	for _, lineup := range lineups {
		var matchups []Matchup
		for _, opponentPlayer := range opponentSkillLevels {
			for _, teamPlayer := range lineup {
				points, err := strconv.ParseFloat(
					strings.Replace(bot.ExcelRows[teamPlayer][opponentPlayer],
						"X", "", 1), 64)
				if err != nil {
					bot.Err = err
					log.Err(bot.Err).Msg("failed to parse float")
					continue
				}
				if points == 0.00 {
					points = 0.01
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
}

// HandlePlayoff for returning max differential expected points lineup from opponent's lineup
func (bot *Data) HandlePlayoff(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	bot.ExcelRows = bot.Excel.GetRows(bot.Excel.WorkBook.Sheets.Sheet[0].Name)
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
}

// HandleCalendar handles the calendar command
func (bot *Data) HandleCalendar(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info().Msg("handling calendar call")
	message := discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				URL:   APACalendarUrl,
				Type:  discordgo.EmbedTypeLink,
				Title: "APA Calendar",
			},
			{
				URL:   TeamCalendarUrl,
				Type:  discordgo.EmbedTypeLink,
				Title: "Team Calendar",
			},
		},
	}
	if m.ChannelID == DevChannelID {
		_, bot.Err = s.ChannelMessageSendComplex(DevChannelID, &message)
	} else {
		_, bot.Err = s.ChannelMessageSendComplex(m.ChannelID, &message)
	}
	if bot.Err != nil {
		log.Err(bot.Err).Msg("failed to post message")
		return
	}
	log.Info().Msgf("calendar posted in Discord channel %s", m.ChannelID)
}
