package communication

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/ytanne/go_nessus/pkg/models"
)

type discordBot struct {
	session        *discordgo.Session
	messageChannel chan models.Message
}

func NewDiscordBot(token string) (*discordBot, error) {
	discordChannel := make(chan models.Message, 10)

	discBot := discordBot{
		messageChannel: discordChannel,
	}

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("Could not start session")
		return nil, err
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(discBot.messageCreate)

	// In this example, we only care about receiving message events.
	discord.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		log.Println("Could not open websocken connection")
		return nil, err
	}

	discBot.session = discord

	return &discBot, nil
}

func (d discordBot) MessageReadChannel() chan models.Message {
	return d.messageChannel
}

func (d discordBot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	d.messageChannel <- models.Message{
		ChannelID: m.ChannelID,
		Msg:       m.Content,
	}
}

func (d discordBot) SendMessage(msg, channelID string) error {
	if channelID == "" {
		return fmt.Errorf("Empty channel ID obtained")
	}

	if _, err := d.session.ChannelMessageSend(channelID, msg); err != nil {
		return err
	}

	return nil
}
