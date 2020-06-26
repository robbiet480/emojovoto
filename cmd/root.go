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
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var botToken string
var guildId string
var channelId string
var approveEmojiId string
var deleteEmojiId string
var abstainEmojiId string
var discordSession *discordgo.Session

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "emojovoto",
	Short: "emojovoto is a easy way to determine what emoji should be deleted from a Discord.",
	Long:  `emojovoto is a easy way to determine what emoji should be deleted from a Discord.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&botToken, "bot-token", "t", "", `bot token to use for authentication. Bot must have read channel, send message, add reactions permissions. Do not prefix with "Bot "`)
	_ = rootCmd.MarkFlagRequired("bot-token")
	rootCmd.PersistentFlags().StringVarP(&guildId, "guild", "g", "", "guild id")
	_ = rootCmd.MarkFlagRequired("guild")
	rootCmd.PersistentFlags().StringVarP(&channelId, "channel", "c", "", "channel id")
	_ = rootCmd.MarkFlagRequired("channel")
	rootCmd.PersistentFlags().StringVarP(&approveEmojiId, "approve-emoji-id", "a", "", "ID of emoji to use to signify approval")
	_ = rootCmd.MarkFlagRequired("approve-emoji-id")
	rootCmd.PersistentFlags().StringVarP(&deleteEmojiId, "delete-emoji-id", "d", "", "ID of emoji to use to signify delete")
	_ = rootCmd.MarkFlagRequired("delete-emoji-id")
	rootCmd.PersistentFlags().StringVarP(&abstainEmojiId, "abstain-emoji-id", "b", "", "ID of emoji to use to signify abstention")
	_ = rootCmd.MarkFlagRequired("abstain-emoji-id")
}

func preRun(cmd *cobra.Command, args []string) {
	if botToken == "" {
		log.Fatalln("Bot token must be set!")
	}

	if guildId == "" {
		log.Fatalln("Guild ID must be set!")
	}

	if channelId == "" {
		log.Fatalln("Channel ID must be set!")
	}

	if approveEmojiId == "" {
		log.Fatalln("Approve emoji ID must be set!")
	}

	if deleteEmojiId == "" {
		log.Fatalln("Delete emoji ID must be set!")
	}

	if abstainEmojiId == "" {
		log.Fatalln("Abstain emoji ID must be set!")
	}

	var connectErr error
	discordSession, connectErr = discordgo.New(fmt.Sprintf("Bot %s", botToken))
	if connectErr != nil {
		log.Fatalln("error creating Discord session", connectErr)
		return
	}
}
