package discord

import "github.com/ytanne/go_nessus/pkg/repository/discord"

type DiscordService interface {
	SendMessage(msg, channelID string) error
}

type discordService struct {
	discordRepo discord.Communicator
}

func NewDiscordService(discordRepo discord.Communicator) DiscordService {
	return discordService{
		discordRepo: discordRepo,
	}
}

func (d discordService) SendMessage(msg, channelID string) error {
	return d.discordRepo.SendMessage(msg, channelID)
}
