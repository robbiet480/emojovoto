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
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var resultsFileName string

// tallyCmd represents the tally command
var tallyCmd = &cobra.Command{
	Use:    "tally",
	Short:  "Tallies up results of ballot",
	Long:   `Tallies up results of ballot. Run this well after you run "ballot". Results are output to results.csv by default.`,
	PreRun: preRun,
	Run: func(cmd *cobra.Command, args []string) {
		results := []TallyResult{}

		messages, messagesErr := getAllMessages(discordSession, channelId)
		if messagesErr != nil {
			log.Fatalln("error getting messages emoji", messagesErr)
			return
		}

		for idx, msg := range messages {
			if !msg.Author.Bot {
				log.Warnf("Skipping message %d of %d as it wasn't sent by a bot", idx+1, len(messages))
				continue
			}

			matches := tallyMsgRegex.FindAllStringSubmatch(msg.Content, -1)

			if len(matches) == 0 {
				log.Warnf("Skipping message %d of %d as it isn't in the expected format", idx+1, len(messages))
				continue
			}

			log.Infof("Processing message %d of %d: %s", idx+1, len(messages), msg.ID)

			match := matches[0]

			resp := TallyResult{
				EmojiID:      match[2],
				EmojiName:    match[1],
				EmojiAddedBy: match[3],
			}

			for _, react := range msg.Reactions {
				if react.Emoji.ID == approveEmojiId { // Approve
					userNames, reactsErr := getUserNames(discordSession, channelId, msg.ID, react.Emoji.APIName())
					if reactsErr != nil {
						log.Errorln("Error when getting users who added thumbs up", reactsErr)
						continue
					}

					resp.KeepCount = len(userNames)
					resp.KeepVoters = userNames
					resp.TotalRespondants = append(resp.TotalRespondants, userNames...)
				} else if react.Emoji.ID == deleteEmojiId { // Delete
					userNames, reactsErr := getUserNames(discordSession, channelId, msg.ID, react.Emoji.APIName())
					if reactsErr != nil {
						log.Errorln("Error when getting users who added thumbs down", reactsErr)
						continue
					}

					resp.DeleteCount = len(userNames)
					resp.DeleteVoters = userNames
					resp.TotalRespondants = append(resp.TotalRespondants, userNames...)
				} else if react.Emoji.ID == abstainEmojiId { // Abstain
					userNames, reactsErr := getUserNames(discordSession, channelId, msg.ID, react.Emoji.APIName())
					if reactsErr != nil {
						log.Errorln("Error when getting users who added abstain", reactsErr)
						continue
					}

					resp.AbstainCount = len(userNames)
					resp.AbstainVoters = userNames
					resp.TotalRespondants = append(resp.TotalRespondants, userNames...)
				}
			}

			resp.TotalRespondantsCount = resp.KeepCount + resp.DeleteCount + resp.AbstainCount

			results = append(results, resp)
		}

		sort.Slice(results, func(i, j int) bool {
			return strings.ToLower(results[i].EmojiName) < strings.ToLower(results[j].EmojiName)
		})

		resultsFile, err := os.OpenFile(resultsFileName, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		if err != nil {
			log.Errorf("Error when opening %s file: %v", resultsFileName, err)
		}
		defer resultsFile.Close()

		if writeErr := gocsv.MarshalFile(&results, resultsFile); writeErr != nil {
			log.Fatalln("Error writing", resultsFileName, writeErr)
		}

		log.Infoln("Saved to", resultsFileName)
	},
}

var tallyMsgRegex = regexp.MustCompile(`Vote here for the <a?:(.*):(.*)> emojo .*, added by (.*)`)

func init() {
	rootCmd.AddCommand(tallyCmd)

	tallyCmd.Flags().StringVarP(&resultsFileName, "file-name", "f", "results.csv", "File name to output results to")
}

func getAllMessages(s *discordgo.Session, channelId string) ([]*discordgo.Message, error) {
	allMessages := []*discordgo.Message{}

	beforeId := "init"
	pageCount := 0

	for beforeId != "" {
		if beforeId == "init" {
			beforeId = ""
		}
		messages, messagesErr := s.ChannelMessages(channelId, 100, beforeId, "", "")
		if messagesErr != nil {
			return nil, messagesErr
		}

		pageCount++

		if len(messages) == 100 {
			beforeId = messages[len(messages)-1].ID
		} else {
			beforeId = ""
		}

		allMessages = append(allMessages, messages...)

		log.Infoln("Got page", pageCount, "of messages")
	}

	log.Infof("Got %d messages", len(allMessages))

	return allMessages, nil
}

func getUserNames(s *discordgo.Session, channelId, msgId, emojiId string) (StringList, error) {
	userNames := StringList{}

	reacts, reactsErr := s.MessageReactions(channelId, msgId, emojiId, 100, "", "")
	if reactsErr != nil {
		return nil, reactsErr
	}

	for _, user := range reacts {
		if user.Bot {
			continue
		}
		userNames = append(userNames, user.Username)
	}

	return userNames, nil
}
