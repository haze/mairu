package main

import (
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

type commandInfo struct {
	Master string
	Slave  []string
}

type eventInfo struct {
	sesh      *discordgo.Session
	message   *discordgo.MessageCreate
	settings  *settings
	config    *Config
	arguments []string
	received  time.Time
}

type settings struct {
	catalyst string
}

var commands map[string]func(eventInfo) (bool, *string)
var eventEmitter *emitter.Emitter
var commandRegex *regexp.Regexp
var gSettings *settings
var gConfig *Config
var gRegistry []*commandInfo

func pushCommand(f func(eventInfo) (bool, *string), name string, aliases ...string) {
	commands[name] = f
	gRegistry = append(gRegistry, &commandInfo{Master: name, Slave: aliases})
}

func registerCommands() {
	pushCommand(PongRoute, "ping", ";p")
	pushCommand(WolframRoute, "?", "?+")
}

func init() {
	t, e := load()
	c, _ := loadConfig()
	if c == nil {
		fmt.Println("Error loading configuration...\nI've made a blank one for you!")
		os.Exit(1)
	} else {
		gConfig = c
	}
	if e != nil {
		panic(e)
	}
	gSettings = t
	eventEmitter = &emitter.Emitter{}
	commandRegex, _ = regexp.Compile("[^\\s\"']+|\"([^\"]*)\"|'([^']*)'")

	commands = make(map[string]func(eventInfo) (bool, *string))
	registerCommands()
}

func searchForCommand(s string) func(eventInfo) (bool, *string) {
	if val, ok := commands[s]; ok {
		return val
	}
	for _, v := range gRegistry {
		for _, z := range v.Slave {
			if strings.ToLower(z) == s {
				return commands[v.Master]
			}
		}
	}
	return nil
}

func route(event *emitter.Event) {
	inf := event.Args[0].(eventInfo)
	message := inf.message.Message.Content
	matches := commandRegex.FindAllString(message, -1)
	inf.arguments = matches
	prefix := strings.ToLower(matches[0])
	f := searchForCommand(prefix)
	if f != nil {
		delete, resp := f(inf)
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
	if gConfig.Token == "" {
		fmt.Println("Missing token in config file!")
		os.Exit(1)
	}

	discord, err := discordgo.New(gConfig.Token)
	if err != nil {
		fmt.Printf("Error: %+v", err)
		os.Exit(1)
	}

	eventEmitter.Use("message", emitter.Void)
	eventEmitter.On("message", func(ev *emitter.Event) {
		go route(ev)
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
	go eventEmitter.Emit("message", eventInfo{sesh: s, config: gConfig, settings: gSettings, received: t, message: m, arguments: []string{}})
}
