package bot

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
	"encoding/json"
	"io"
	"os"
	"sync"
)

// Schedule represents a single match in the team's schedule
type Schedule struct {
	Date     string `json:"Date"`
	Opponent string `json:"Opponent"`
}

// TeamSchedule represents the schedule for a single team
type TeamSchedule struct {
	Team     string     `json:"Team"`
	Schedule []Schedule `json:"Schedule"`
}

// Schedules represents the overall structure of the JSON file
type Schedules struct {
	Schedules []TeamSchedule `json:"Schedules"`
}

// LoadSchedules loads the schedules from a JSON file
func LoadSchedules(schedules *Schedules, path string) error {
	file, err := os.Open(path)
	if err != nil {
			return err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
			return err
	}

	err = json.Unmarshal(bytes, schedules)
	if err != nil {
			return err
	}

	return nil
}

// ScheduleGameDay for posting game day message
func (bot *Data) ScheduleGameDay(s *discordgo.Session, m *discordgo.MessageCreate, teamName string) {
	log.Info().Msg("handling game day post creation")

	var schedules Schedules
	err := LoadSchedules(&schedules, "data/schedules/Spring2025Schedule.json")
	if err != nil {
			bot.Err = err
			log.Err(bot.Err).Msg("failed to load schedules")
			return
	}

	var wg sync.WaitGroup

	for _, teamSchedule := range schedules.Schedules {
			if teamSchedule.Team != teamName {
					continue
			}
			for matchIndex, match := range teamSchedule.Schedule {
					matchDate, err := time.Parse("01/02/2006", match.Date)
					if err != nil {
							log.Err(err).Msgf("failed to parse match date: %s", match.Date)
							continue
					}
					if matchDate.Before(time.Now()) {
							continue
					}
					// TODO: revert to -1 and 9
					// Schedule the post for the day before the match at 9am EST
					postTime := matchDate.AddDate(0, 0, -2).Add(15 * time.Hour)
					loc, err := time.LoadLocation("America/New_York")
					if err != nil {
							loc = time.UTC
							log.Err(err).Msg("failed to load timezone, using UTC as EST+5")
					}
					postTime = postTime.In(loc)
					log.Info().Msgf("postTime: %s, current time: %s, until: %v", postTime, time.Now(), time.Until(postTime))
					duration := time.Until(postTime)
					if duration <= 0 {
							log.Warn().Msgf("Post time %s is in the past or immediate; skipping scheduling", postTime)
							continue
					}


					
					wg.Add(1)
					time.AfterFunc(time.Until(postTime), func() {
							defer wg.Done()
							var team Team
							if teamSchedule.Team == WookieMistakes.Name {
									team = WookieMistakes
							} else if teamSchedule.Team == SafetyDance.Name {
									team = SafetyDance
							} else {
									bot.Err = errors.New("invalid team name")
									log.Err(bot.Err).Msgf("failed to create game day post for team: %s", teamSchedule.Team)
									return
							}

							var customMessage string

							var playbacksMessage string
							if  matchIndex >= len(teamSchedule.Schedule)-2 {
								playbacksMessage = " with no playbacks"
						}
							message := discordgo.MessageSend{
									Content: fmt.Sprintf(
											"@everyone Attendance time <a:abongoblob:1324456047661813851> This week %s plays %s%s <a:Toothless:1324460455623655535>\n"+
													ReactionRequest+customMessage, team.Name, match.Opponent, playbacksMessage,
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
									for _, t := range teammate.Teams {
											if t.Name == team.Name {
													numspaces = longestName + 1 - len(teammate.LastName)
													var skillLevel int
													switch team.Name {
													case WookieMistakes.Name:
															skillLevel = teammate.SkillLevel.Eight
													case SafetyDance.Name:
															skillLevel = teammate.SkillLevel.Nine
													default:
															bot.Err = errors.New("invalid team name")
															log.Err(bot.Err).Msgf("failed to create game day post for team: %s", team.Name)
															return
													}
													message.Content += fmt.Sprintf("|%s| %s%s|⬛|⬛|⬛|⬛|\n", intToEmoji(skillLevel), teammate.LastName, strings.Repeat(" ", numspaces))
											}
									}
							}
							message.Content += fmt.Sprintf("+➖+%s+➖+➖+➖+➖+\n```", strings.Repeat("-", longestName+2))
							message.Content += "Eligible Lineups:"
							// TODO: revert to team.GameNightChannelID
							postedMessage, err := s.ChannelMessageSendComplex(DevChannelID, &message)
							if err != nil {
									bot.Err = err
									log.Err(bot.Err).Msg("failed to post message")
									return
							}
							emotes := []string{"👍", "⏳", "👎", "❓"}
							for _, emote := range emotes {
									err := s.MessageReactionAdd(postedMessage.ChannelID, postedMessage.ID, emote)
									if err != nil {
											log.Err(err).Msgf("failed to add reaction %s to message %s", emote, postedMessage.ID)
									}
							}
							log.Info().Msgf("game day %s vs %s posted or scheduled to Discord channel %s", team.Name, match.Opponent, m.ChannelID)
					})
			}
	}
	wg.Wait()
}