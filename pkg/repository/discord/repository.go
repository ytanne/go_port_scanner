package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/ytanne/go_nessus/pkg/models"
)

var DiscordChannel chan models.Message

type discordBot struct {
	session *discordgo.Session
}

type Communicator interface {
	SendMessage(msg, channelID string) error
}

func NewDiscordBot(token string) (Communicator, error) {
	DiscordChannel = make(chan models.Message, 10)

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("Could not start session")
		return nil, err
	}
	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	discord.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		log.Println("Could not open websocken connection")
		return nil, err
	}
	return discordBot{
		session: discord,
	}, nil
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	DiscordChannel <- models.Message{
		ChannelID: m.ChannelID,
		Msg:       m.Content,
	}
}

func (d discordBot) SendMessage(msg, channelID string) error {
	if _, err := d.session.ChannelMessageSend(channelID, msg); err != nil {
		log.Println("Could not send message to channel", channelID, ". Error:", err)
		return err
	}
	return nil
}
