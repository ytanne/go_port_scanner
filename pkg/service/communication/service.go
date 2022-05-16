package communication

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/ytanne/go_port_scanner/pkg/models"
	"os"
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
		return nil, fmt.Errorf("could not start new discord bot: %w", err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(discBot.messageCreate)

	// In this example, we only care about receiving message events.
	discord.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		return nil, fmt.Errorf("could not open websocket connection: %w", err)
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
		return fmt.Errorf("empty channel ID obtained")
	}

	if _, err := d.session.ChannelMessageSend(channelID, msg); err != nil {
		return fmt.Errorf("could not send msg to channel: %w", err)
	}

	return nil
}

func (d discordBot) SendFile(channelID, fileName string) error {
	fd, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("could not open %s. Error: %w", fileName, err)
	}

	defer fd.Close()

	_, err = d.session.ChannelFileSend(channelID, fileName, fd)
	if err != nil {
		return fmt.Errorf("could not send file %s. Error: %w", fileName, err)
	}

	return nil
}
