/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
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
	msgContent := fmt.Sprintf("Vote here for the %s emojo (name: `%s`, id: `%s`), added by %s", emojo.MessageFormat(), emojo.Name, emojo.ID, emojo.User.Username)
	msg := &discordgo.MessageSend{
		Content: msgContent,
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
