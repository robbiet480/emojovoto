/*
Copyright © 2020 Robbie Trencheny

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var approveEmojiFormatted string
var deleteEmojiFormatted string
var abstainEmojiFormatted string

// ballotCmd represents the ballot command
var ballotCmd = &cobra.Command{
	Use:    "ballot",
	Short:  "Create the voting ballot in the channel",
	Long:   `Create the voting ballot in the channel. Will get a list of all emoji in the guild, send a message in the specified channel and then add the three reacts on it.`,
	PreRun: preRun,
	Run: func(cmd *cobra.Command, args []string) {
		// Bot must open gateway connection at least once before being able to send messages
		if openErr := discordSession.Open(); openErr != nil {
			log.Fatalln("error opening connection", openErr)
			return
		}

		if closeErr := discordSession.Close(); closeErr != nil {
			log.Fatalln("error closing connection", closeErr)
			return
		}

		emojiList, emojiListErr := GuildEmoji(discordSession, guildId)
		if emojiListErr != nil {
			log.Fatalln("error listing emoji", emojiListErr)
			return
		}

		for _, emoji := range emojiList {
			if emoji.ID == approveEmojiId {
				approveEmojiFormatted = emoji.APIName()
			} else if emoji.ID == deleteEmojiId {
				deleteEmojiFormatted = emoji.APIName()
			} else if emoji.ID == abstainEmojiId {
				abstainEmojiFormatted = emoji.APIName()
			}
		}

		sort.Slice(emojiList, func(i, j int) bool {
			return strings.ToLower(emojiList[i].Name) < strings.ToLower(emojiList[j].Name)
		})

		for _, emojo := range emojiList {
			if emojoSendErr := createEmojoMessage(discordSession, channelId, emojo); emojoSendErr != nil {
				log.Errorln("Error writing emojo message", emojoSendErr)
			}

			time.Sleep(3 * time.Second)
		}
	},
}

func init() {
	rootCmd.AddCommand(ballotCmd)
}

func createEmojoMessage(s *discordgo.Session, channelId string, emojo *Emoji) error {
	msg := &discordgo.MessageSend{
		Content: emojo.MessageFormat(),
		Embed: &discordgo.MessageEmbed{
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Name",
					Value:  emojo.Name,
					Inline: true,
				},
				{
					Name:   "ID",
					Value:  emojo.ID,
					Inline: true,
				},
				{
					Name:   "Added By",
					Value:  emojo.User.Username,
					Inline: true,
				},
			},
		},
	}

	sentMsg, msgErr := s.ChannelMessageSendComplex(channelId, msg)
	if msgErr != nil {
		return msgErr
	}

	// approve, delete, abstain
	for _, react := range []string{approveEmojiFormatted, deleteEmojiFormatted, abstainEmojiFormatted} {
		if reactErr := s.MessageReactionAdd(channelId, sentMsg.ID, react); reactErr != nil {
			log.Errorf("Error when adding %s emoji to %s msg: %v", react, sentMsg.ID, reactErr)
		}
	}

	return nil
}

// GuildEmoji returns all emoji
// guildID : The ID of a Guild.
func GuildEmoji(s *discordgo.Session, guildID string) (emoji []*Emoji, err error) {

	body, err := s.RequestWithBucketID("GET", discordgo.EndpointGuildEmojis(guildID), nil, discordgo.EndpointGuildEmojis(guildID))
	if err != nil {
		return
	}

	if err := json.Unmarshal(body, &emoji); err != nil {
		return nil, discordgo.ErrJSONUnmarshal
	}
	return
}
