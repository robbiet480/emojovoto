# emojovoto

emojovoto is a easy way to determine what emoji should be deleted from a Discord

## Process

1. Create a Discord bot user
2. Add it to your server with at least send messages, add reactions and read message history permissions.
3. `go build -o emojovoto .` in the root directory
4. Run `emojovoto ballot` to output the "ballot". Ballot looks like one message per emoji on your server with 3 reactions added to it: approve, delete and abstain.
5. Wait some time for enough server members to vote.
6. Run `emojovoto tally` to output the vote results to a CSV, `results.csv` by default.
7. Based on your own criteria, delete emojis as you wish.

## Help

```
emojovoto is a easy way to determine what emoji should be deleted from a Discord.

Usage:
  emojovoto [command]

Available Commands:
  ballot      Create the voting ballot in the channel
  help        Help about any command
  tally       Tallies up results of ballot

Flags:
  -b, --abstain-emoji-id string   ID of emoji to use to signify abstention (default "726149181458088016")
  -a, --approve-emoji-id string   ID of emoji to use to signify approval (default "726149138651021413")
  -t, --bot-token string          bot token to use for authentication. Bot must have read channel, send message, add reactions permissions. Do not prefix with "Bot "
  -c, --channel string            channel id
  -d, --delete-emoji-id string    ID of emoji to use to signify delete (default "726149161195667506")
  -g, --guild string              guild id
  -h, --help                      help for emojovoto

Use "emojovoto [command] --help" for more information about a command.
```

## TODO

- Add auto delete emoji after tally

## LICENSE

MIT
