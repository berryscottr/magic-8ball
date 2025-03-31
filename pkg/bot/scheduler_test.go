package bot

import (
    "errors"
    "fmt"
    "os"
    "strings"
    "sync"
    "testing"
    "time"

    "github.com/bwmarrin/discordgo"
    "github.com/stretchr/testify/assert"
    "github.com/rs/zerolog/log"
)

// TestScheduleGameDay for posting game day message after 5 seconds
func TestScheduleGameDay(t *testing.T) {
    assertion := assert.New(t)
    botToken := os.Getenv("BOT_TOKEN")
    if botToken == "" {
        t.Fatal("BOT_TOKEN not set")
    }
    data := Data{Token: Token{Discord: botToken}}
    data.SetDir()

    session, err := discordgo.New("Bot " + data.Token.Discord)
    assertion.NoError(err, "failed to create Discord session")

    // Open a websocket connection to Discord
    err = session.Open()
    assertion.NoError(err, "failed to open Discord session")
    defer session.Close()

    var schedules Schedules
    err = LoadSchedules(&schedules, "../../data/schedules/Spring2025Schedule.json")
    assertion.NoError(err, "failed to load schedules")

    var wg sync.WaitGroup

    for _, teamSchedule := range schedules.Schedules {
        if teamSchedule.Team == "Wookie Mistakes" || teamSchedule.Team == "Safety Dance" {
            if len(teamSchedule.Schedule) == 0 {
                t.Logf("no matches found for team: %s", teamSchedule.Team)
                continue
            }

            nextMatch := teamSchedule.Schedule[0]
            wg.Add(1)
            time.AfterFunc(5*time.Second, func() {
                defer wg.Done()

                var team Team
                if teamSchedule.Team == WookieMistakes.Name {
                    team = WookieMistakes
                } else if teamSchedule.Team == SafetyDance.Name {
                    team = SafetyDance
                } else {
                    data.Err = errors.New("invalid team name")
                    t.Errorf("failed to create game day post for team: %s", teamSchedule.Team)
                    return
                }

                var opponentTeam string
                for i, name := range team.DivisionTeamNames {
                    for _, junk := range []string{"'", "-", "8", "9"} {
                        name = strings.Replace(name, junk, "", -1)
                    }
                    name = strings.TrimRight(name, " ")
                    if strings.Contains(strings.ToLower(nextMatch.Opponent), strings.ToLower(name)) {
                        opponentTeam = team.DivisionTeamNames[i]
                    }
                }

                message := discordgo.MessageSend{
                    Content: fmt.Sprintf(
                        "@everyone Attendance time <a:abongoblob:1324456047661813851> This week %s plays %s <a:Toothless:1324460455623655535>\n"+
                            ReactionRequest, team.Name, opponentTeam,
                    ),
                }
                message.Content += "\n```\n"
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
                message.Content += "+🎱+---Name---+👍+⏳+👎+❓+\n"
                var numspaces int
                for _, teammate := range Teammates {
                    for _, tt := range teammate.Teams {
                        if tt.Name == team.Name {
                            numspaces = longestName + 1 - len(teammate.LastName)
                            var skillLevel int
                            switch team.Name {
                            case WookieMistakes.Name:
                                skillLevel = teammate.SkillLevel.Eight
                            case SafetyDance.Name:
                                skillLevel = teammate.SkillLevel.Nine
                            default:
                                data.Err = errors.New("invalid team name")
                                t.Errorf("failed to create game day post for team: %s", team.Name)
                                return
                            }
                            message.Content += fmt.Sprintf("|%s| %s%s|⬛|⬛|⬛|⬛|\n", intToEmoji(skillLevel), teammate.LastName, strings.Repeat(" ", numspaces))
                        }
                    }
                }
                message.Content += fmt.Sprintf("+➖+%s+➖+➖+➖+➖+\n```", strings.Repeat("-", longestName+2))
                message.Content += "Eligible Lineups:"
                postedMessage, err := session.ChannelMessageSendComplex(DevChannelID, &message)
                if err != nil {
                    data.Err = err
                    log.Err(data.Err).Msg("failed to post message")
                    t.Errorf("failed to post message")
                    return
                }
                emotes := []string{"👍", "⏳", "👎", "❓"}
                for _, emote := range emotes {
                    err := session.MessageReactionAdd(postedMessage.ChannelID, postedMessage.ID, emote)
                    if err != nil {
                        t.Errorf("failed to add reaction %s to message %s", emote, postedMessage.ID)
                    }
                }
                t.Logf("test game day %s vs %s posted to Discord channel %s", team.Name, opponentTeam, DevChannelID)
            })
        }
    }

    wg.Wait()
}