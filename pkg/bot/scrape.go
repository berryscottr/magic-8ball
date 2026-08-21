package bot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
    "strconv"
	"strings"
	"sync"

	"github.com/berryscottr/magic-8ball/pkg/scrapaer"
	"github.com/berryscottr/magic-8ball/pkg/util"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// HandleScrape calls upon Scrapaer package to scrape APA site using access token
func (bot *Data) HandleScrape(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info().Msg("handling scrape")

	bot.Token.APA, bot.Err = util.RunPythonScript("scripts/scrape/main.py")
    bot.Token.APA = strings.TrimSpace(bot.Token.APA)
	teamIDs := []int{WookieMistakes8.ID, WookieMistakes9.ID}
	memberIDs := []int{}
    
	teamIDs = overrideIntListFlag(m.Content, "-teams", teamIDs)
	memberIDs = overrideIntListFlag(m.Content, "-members", memberIDs)

	playersScraped := scrapeTeams(bot.Token.APA, teamIDs, memberIDs)

	var skillEvals string
    skillEvals, bot.Err = util.RunPythonScript("scripts/equalizer/main.py")

    scrapedNames := make(map[string]struct{})
    for _, p := range playersScraped {
        scrapedNames[p] = struct{}{}
    }

    lines := strings.Split(skillEvals, "\n")
    var filtered []string

    for _, line := range lines {
        for name := range scrapedNames {
            if strings.Contains(line, name) {
                filtered = append(filtered, line)
                break
            }
        }
    }

    skillEvals = strings.Join(filtered, "\n")

	message := discordgo.MessageSend{}
	message.Content = skillEvals
	if m.ChannelID == DevChannelID {
		_, bot.Err = s.ChannelMessageSendComplex(DevChannelID, &message)
	} else if m.ChannelID == TestChannelID {
		_, bot.Err = s.ChannelMessageSendComplex(TestChannelID, &message)
	} else {
		_, bot.Err = s.ChannelMessageSendComplex(StrategyChannelID, &message)
	}
	if bot.Err != nil {
		log.Err(bot.Err).Msg("failed to post message")
		return
	}
	log.Info().Msgf("%v scraped content posted to Discord channel %s", "", m.ChannelID)
}

// overrideIntListFlag parses flags like:  -teams=1,2,3
// If found in text, return those values. Otherwise return original.
func overrideIntListFlag(text, flag string, original []int) []int {
	// Example regex: -teams={1,2,3}
	pattern := fmt.Sprintf(`%s=([^}]*)`, regexp.QuoteMeta(flag))
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(text)
	if len(matches) < 2 {
		return original // no override
	}

	// matches[1] = "1,2,3"
	raw := matches[1]
	rawParts := strings.Split(raw, ",")

	var results []int
	for _, part := range rawParts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		val, err := strconv.Atoi(part)
		if err != nil {
			continue // silently skip bad values
		}
		results = append(results, val)
	}

	if len(results) == 0 {
		return original // nothing usable
	}
	return results
}


var fileMutex sync.Mutex
var mu sync.Mutex

func scrapeTeams(authToken string, teamIds []int, memberIds []int) []string {
    headers := map[string]string{
        "user-agent":    "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.104 Safari/537.36",
        "referer":       "https://league.poolplayers.com",
        "Authorization": authToken,
        "authority":     "gql.poolplayers.com",
    }

    graphqlURL := "https://gql.poolplayers.com/graphql"
    client := scrapaer.NewGraphQLClient(graphqlURL, headers)
    api := scrapaer.NewPoolPlayersAPI(client)

    var wg sync.WaitGroup
    if len(teamIds) > 0 {
        for _, teamID := range teamIds {
            wg.Add(1)
            go func(teamID int) {
                defer wg.Done()
                teamRosterResult, err := api.GetTeamRoster(teamID)
                if err != nil {
										log.Err(err).Msgf("Error fetching team roster for team %d", teamID)
                    return
                }
                for _, teamMember := range teamRosterResult.Team.Roster {
                    memberID := teamMember.Member.ID
                    if memberID != 0 {
                        mu.Lock()
                        memberIds = appendUniqueInt(memberIds, memberID)
                        mu.Unlock()
                    }
                }
            }(teamID)
        }
        wg.Wait()
    }
    var playersScraped []string
    for _, memberID := range memberIds {
        wg.Add(1)
        go func(memberID int) {
            defer wg.Done()
            playerName := runPlayerReport(api, memberID)
            if playerName != "" {
                mu.Lock()
                playersScraped = append(playersScraped, playerName)
                mu.Unlock()
            }
        }(memberID)
    }

    wg.Wait()
    return playersScraped
}

func runPlayerReport(api *scrapaer.PoolPlayersAPI, memberID int) string {
    var player scrapaer.Player
    player.MemberID = memberID
    memberStatsHeaderResult, err := api.GetMemberStatsHeader(player.MemberID)
    if err != nil {
				log.Err(err).Msg("Invalid Member")
        return ""
    }
    var nonAlphaRegex = regexp.MustCompile("[^a-zA-Z]+")
    player.Name = fmt.Sprintf("%s %s", nonAlphaRegex.ReplaceAllString(memberStatsHeaderResult.Member.FirstName, ""), nonAlphaRegex.ReplaceAllString(memberStatsHeaderResult.Member.LastName, ""))
    player.PlayerID = memberStatsHeaderResult.Member.Aliases[0].ID
    playerResult, err := api.GetPlayerTeams(player.PlayerID)
    if err != nil {
        log.Err(err).Msg("Invalid Player")
        return ""
    }
    for _, currentTeam := range playerResult.Alias.CurrentTeams {
        if currentTeam.Team.ID != 0 {
            player.TeamIDs = append(player.TeamIDs, currentTeam.Team.ID)
        }
    }
    for _, pastTeam := range playerResult.Alias.PastTeams {
        if pastTeam.Team.ID != 0 {
            player.TeamIDs = append(player.TeamIDs, pastTeam.Team.ID)
        }
    }

    var wg sync.WaitGroup
    matchesCh := make(chan scrapaer.PlayerMatch)

    for _, teamID := range player.TeamIDs {
        wg.Add(1)
        go func(teamID int) {
            defer wg.Done()
            teamRosterResult, err := api.GetTeamRoster(teamID)
            if err != nil {
                log.Err(err).Msgf("Error fetching team roster for team %d", teamID)
                return
            }
            var teamPlayerID int
            for _, teamMember := range teamRosterResult.Team.Roster {
                if teamMember.Member.ID == player.MemberID {
                    teamPlayerID = teamMember.ID
                }
            }
            teamScheduleResult, err := api.GetTeamSchedule(teamID)
            if err != nil {
								log.Err(err).Msgf("Error fetching team schedule for team %d", teamID)
                return
            }

            var matchWg sync.WaitGroup
            for _, teamMatch := range teamScheduleResult.Team.Matches {
                teamMatchID, ok := teamMatch.ID.(float64)
                if ok && teamMatchID != 0 {
                    matchWg.Add(1)
                    go func(teamMatchID int) {
                        defer matchWg.Done()
                        matchResult, err := api.GetMatch(teamMatchID)
                        if err != nil {
                            log.Err(err).Msgf("Error fetching match %d", teamMatchID)
                            return
                        }
                        if matchResult.Match.Type == "EIGHT" || matchResult.Match.Type == "NINE" {
                            processMatch(matchResult.Match.Type, matchResult, teamPlayerID, matchesCh)
                        }
                    }(int(teamMatchID))                    
                }
            }
            matchWg.Wait()
        }(teamID)
    }

    go func() {
        wg.Wait()
        close(matchesCh)
    }()

    for score := range matchesCh {
        player.Matches = append(player.Matches, score)
    }

    err = matches2json(player)
    if err != nil {
        log.Err(err).Msg("failed to save Matched to JSON")
    }
    return player.Name
}

func getPlayerScores(results []scrapaer.MatchResult, teamPlayerID int) []scrapaer.PlayerMatch {
    var playerScores []scrapaer.PlayerMatch
    for _, team := range results {
        for _, score := range team.Scores {
            if score.Player.ID == teamPlayerID {
                playerScores = append(playerScores, score)
            }
        }
    }
    return playerScores
}

func getOpponentStats(results []scrapaer.MatchResult, matchPosition, teamPlayerID int) (int, int, float64, float64) {
    var opponentSkillLevel, opponentDefensiveShots int
    var eightBallLosses, eightBallMatchPointsLost float64

    for _, team := range results {
        for _, score := range team.Scores {
            if score.MatchPositionNumber == matchPosition && score.Player.ID != teamPlayerID {
                opponentSkillLevel = score.SkillLevel
                opponentDefensiveShots = score.DefensiveShots
                if value, ok := score.EightBallWins.(float64); ok {
                    eightBallLosses = value
                }
                if value, ok := score.EightBallMatchPointsEarned.(float64); ok {
                    eightBallMatchPointsLost = value
                }
            }
        }
    }
    return opponentSkillLevel, opponentDefensiveShots, eightBallLosses, eightBallMatchPointsLost
}

func processMatch(matchType string, matchResult scrapaer.MatchPage, teamPlayerID int, matchesCh chan<- scrapaer.PlayerMatch) {
    tableSize, ok := scrapaer.TableSizes[matchResult.Match.Location.ID]
    if !ok {
				log.Warn().Msgf("No table size found for location %s:%d, using default 8", matchResult.Match.Location.Name, matchResult.Match.Location.ID)
        tableSize = 8
    }
    playerScores := getPlayerScores(matchResult.Match.Results, teamPlayerID)

    for i := range playerScores {
        matchPosition := playerScores[i].MatchPositionNumber
        opponentSkillLevel, opponentDefensiveShots, eightBallLosses, eightBallMatchPointsLost := getOpponentStats(matchResult.Match.Results, matchPosition, teamPlayerID)

        playerScores[i].OpponentSkillLevel = opponentSkillLevel
        playerScores[i].OpponentDefensiveShots = opponentDefensiveShots

        if matchType == "EIGHT" {
            playerScores[i].EightBallLosses = eightBallLosses
            playerScores[i].EightBallMatchPointsLost = eightBallMatchPointsLost
        }

        playerScores[i].TableSize = tableSize

        matchesCh <- playerScores[i]
    }
}

func matches2json(player scrapaer.Player) error {
    jsonData, err := json.MarshalIndent(player.Matches, "", "  ")
    if err != nil {
        log.Error().Err(err).Msg("Error marshalling JSON")
        return err
    }

    dir := fmt.Sprintf("data/players/%s", player.Name)
    filePath := filepath.Join(dir, "matches.json")

    fileMutex.Lock()
    defer fileMutex.Unlock()

    if err := os.MkdirAll(dir, os.ModePerm); err != nil {
        log.Error().Err(err).Str("dir", dir).Msg("Error creating directory")
        return err
    }

    file, err := os.Create(filePath)
    if err != nil {
        log.Error().Err(err).Str("path", filePath).Msg("Error creating file")
        return err
    }
    defer file.Close()

    _, err = file.Write(jsonData)
    if err != nil {
        log.Error().Err(err).Str("path", filePath).Msg("Error writing JSON")
        return err
    }

    log.Info().Str("path", filePath).Msg("JSON data saved successfully")
    return nil
}


func appendUniqueInt(slice []int, elems ...int) []int {
	seen := make(map[int]bool)
	for _, elem := range slice {
		seen[elem] = true
	}

	for _, elem := range elems {
		if !seen[elem] {
			slice = append(slice, elem)
			seen[elem] = true
		}
	}
	return slice
}

