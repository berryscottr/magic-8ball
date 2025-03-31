package bot

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// HandleScrape calls upon Scrapaer package to scrape APA site using access token
func (bot *Data) HandleScrape(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info().Msg("handling scrape")
	re := regexp.MustCompile(`"(.*?)"`)
	var content string
	contentSlice := strings.Split(m.Content, "!")
	for _, com := range contentSlice {
		if strings.Contains(com, "scrape") {
			content = com
		}
	}
	bot.Token.APA = re.FindString(content)
	var memberIDs string
	for _, teammate := range Teammates {
		if len(teammate.Teams) > 0 {
			if memberIDs == "" {
				memberIDs = fmt.Sprintf("%d", teammate.ID.Member)
			}
			memberIDs = fmt.Sprintf("%s,%d", memberIDs, teammate.ID.Member)
		}
	}
	// TODO: execute call of scrapaer endpoint for current functionality
	// TODO: add an arg to supply team ID to scrapaer endpoint (scrapaer will need to be updated)
	message := discordgo.MessageSend{}
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
