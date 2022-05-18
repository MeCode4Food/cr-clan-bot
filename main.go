package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mecode4food/cr-clan-bot/pkg/clashroyale"
	"github.com/mecode4food/cr-clan-bot/pkg/config"
)

var (
	d *discordgo.Session

	aid = config.Viper().GetString("bot.app_id")     // app id
	cid = config.Viper().GetString("bot.channel_id") // channel id
	gid = config.Viper().GetString("bot.guild_id")   // guild (server) id
	t   = config.Viper().GetString("bot.token")      // bot token

	commands = map[string]*discordgo.ApplicationCommand{
		"clan": {
			Name:        "clan",
			Description: "Get clan details",
		},
		"echo": {
			Name:        "echo",
			Description: "Echo a message to the channel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "message",
					Description: "Message to echo",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		"help": {
			Name:        "help",
			Description: "List all commands",
		},
		"inactive": {
			Name:        "inactive",
			Description: "Get list of inactive members",
		},
	}
	handlers = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"clan":     clan,
		"echo":     messageEcho,
		"help":     help,
		"inactive": inactive,
	}
)

func main() {
	var err error
	d, err = discordgo.New("Bot " + t)
	if err != nil {
		log.Fatal(err)
	}
	d.AddHandler(botReady)

	for _, cmd := range commands {
		_, err := d.ApplicationCommandCreate(aid, gid, cmd)
		if err != nil {
			log.Fatal(err)
		}
	}

	d.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		d := i.ApplicationCommandData()
		if v, ok := handlers[d.Name]; ok {
			v(s, i)
		}
	})

	err = d.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer d.Close()

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-sc
	log.Println("Bot is now closing")
	// remove the commands
	cc, err := d.ApplicationCommands(aid, gid)
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range cc {
		d.ApplicationCommandDelete(aid, gid, c.ID)
	}
}

func botReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("Bot is ready")
}

func clan(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var m string
	id := config.Viper().GetString("clash_royale.clan_id")
	clan, err := clashroyale.Clan(id)
	if err != nil {
		log.Println(err)
		m = "Error getting clan info"
	} else {
		m = fmt.Sprintf("Clan: %s\nTag: %s\nMembers: %d\n", clan.Name, clan.Tag, clan.Members)
	}
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: m,
		},
	})
	if err != nil {
		log.Println(err)
	}
}

func inactive(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var msg string
	msg = "Inactive members in Electro Shack:\n"
	id := config.Viper().GetString("clash_royale.clan_id")
	mm, err := clashroyale.ClanMembers(id)

	if err != nil {
		log.Println(err)
		msg = "Error getting clan members info"
	} else {
		sort.Slice(mm, func(i2, j int) bool {
			return mm[i2].LastSeen().Before(mm[j].LastSeen())
		})
		for _, m := range mm {
			if time.Since(m.LastSeen()).Hours() > 24 {
				msg += fmt.Sprintf("* **%s**: %s Trophies: %d Last Seen: **%.2f Days ago**\n", m.Name, m.Tag, m.Trophies, time.Since(m.LastSeen()).Hours()/24)
			}
		}
	}
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	if err != nil {
		log.Println(err)
	}
}

func messageEcho(s *discordgo.Session, i *discordgo.InteractionCreate) {
	d := i.ApplicationCommandData()
	m := fmt.Sprintf("%s said: %s", i.Member.User.Username, d.Options[0].Value)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: m,
		},
	})
	if err != nil {
		log.Println(err)
	}
}

func help(s *discordgo.Session, i *discordgo.InteractionCreate) {
	m := "List of commands:\n```"
	for _, cmd := range commands {
		m += fmt.Sprintf("%s: %s\n", cmd.Name, cmd.Description)
	}
	extra := "Welcome to the Electro Shack clan chat!\nIf you have suggestions (or you want to help add commands to this bot), feel free to ping chick!"
	m += "```\n" + extra

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: m,
		},
	})
	if err != nil {
		log.Println(err)
	}
}
