package cmd

import (
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Emoji struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Roles         []string        `json:"roles"`
	User          *discordgo.User `json:"user"`
	RequireColons bool            `json:"require_colons"`
	Managed       bool            `json:"managed"`
	Animated      bool            `json:"animated"`
	Available     bool            `json:"available"`
}

// MessageFormat returns a correctly formatted Emoji for use in Message content and embeds
func (e *Emoji) MessageFormat() string {
	if e.ID != "" && e.Name != "" {
		if e.Animated {
			return "<a:" + e.APIName() + ">"
		}

		return "<:" + e.APIName() + ">"
	}

	return e.APIName()
}

// APIName returns an correctly formatted API name for use in the MessageReactions endpoints.
func (e *Emoji) APIName() string {
	if e.ID != "" && e.Name != "" {
		return e.Name + ":" + e.ID
	}
	if e.Name != "" {
		return e.Name
	}
	return e.ID
}

type StringList []string

func (s StringList) MarshalCSV() (string, error) {
	orig := []string{}

	for _, each := range s {
		temp := string(each)
		orig = append(orig, temp)
	}

	sort.Slice(orig, func(i, j int) bool { return strings.ToLower(orig[i]) < strings.ToLower(orig[j]) })

	return strings.Join(orig, ", "), nil
}

type TallyResult struct {
	EmojiID      string
	EmojiName    string
	EmojiAddedBy string

	KeepCount    int
	DeleteCount  int
	AbstainCount int

	TotalRespondantsCount int
	TotalRespondants      StringList

	KeepVoters    StringList
	DeleteVoters  StringList
	AbstainVoters StringList
}
