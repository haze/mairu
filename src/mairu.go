package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/olebedev/emitter"
)

type eventInfo struct {
	sesh      *discordgo.Session
	message   *discordgo.MessageCreate
	arguments []string
	received  time.Time
}

type settings struct {
	catalyst string
}

var commands map[string]func(eventInfo) (bool, *string)
var token string
var eventEmitter *emitter.Emitter
var commandRegex *regexp.Regexp
var gSettings *settings

func registerCommands() {
	commands["ping"] = PongRoute
}

func init() {
	t, e := load()
	if e != nil {
		panic(e)
	}
	gSettings = t
	eventEmitter = &emitter.Emitter{}
	commandRegex, _ = regexp.Compile("[^\\s\"']+|\"([^\"]*)\"|'([^']*)'")

	commands = make(map[string]func(eventInfo) (bool, *string))
	registerCommands()

	flag.StringVar(&token, "t", "", "user token")
	flag.Parse()
}

func route(event *emitter.Event) {
	inf := event.Args[0].(eventInfo)
	message := inf.message.Message.Content
	matches := commandRegex.FindAllString(message, -1)
	prefix := strings.ToLower(matches[0])
	if val, ok := commands[prefix]; ok {
		delete, resp := val(inf)
		if delete {
			inf.sesh.ChannelMessageDelete(inf.message.ChannelID, inf.message.ID)
			if resp != nil {
				inf.sesh.ChannelMessageSend(inf.message.ChannelID, *resp)
			}
		} else {
			inf.sesh.ChannelMessageEdit(inf.message.ChannelID, inf.message.ID, *resp)
		}
	}
}

func main() {

	if token == "" {
		fmt.Println("User token not provided; usage: mairu -t <token>")
		os.Exit(1)
	}

	discord, err := discordgo.New(token)
	if err != nil {
		fmt.Printf("Error: %+v", err)
		os.Exit(1)
	}

	eventEmitter.Use("message", emitter.Void)
	eventEmitter.On("message", func(ev *emitter.Event) {
		route(ev)
	})

	// Register our shizz.

	discord.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
		os.Exit(1)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Mairu is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	t := time.Now()
	if m.Author.ID != s.State.User.ID {
		return
	}
	// debug print
	go eventEmitter.Emit("message", eventInfo{sesh: s, received: t, message: m, arguments: []string{}})
}
