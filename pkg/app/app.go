package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ytanne/go_nessus/pkg/models"
	"github.com/ytanne/go_nessus/pkg/repository/discord"
	"github.com/ytanne/go_nessus/pkg/service/nmap"
	"github.com/ytanne/go_nessus/pkg/service/sqlite"
)

type App struct {
	communicator discord.Communicator
	storage      sqlite.DBKeeper
	portScanner  nmap.NmapScanner
	channelType  map[int]string
}

func NewApp(communicator discord.Communicator, storage sqlite.DBKeeper, portScanner nmap.NmapScanner) *App {
	return &App{
		communicator: communicator,
		storage:      storage,
		portScanner:  portScanner,
		channelType:  make(map[int]string),
	}
}

func (c *App) SetUpChannels(arpChannelID, psChannelID, wpsChannelID string) {
	if arpChannelID != "" {
		c.channelType[models.ARP] = arpChannelID
	}
	if psChannelID != "" {
		c.channelType[models.PS] = psChannelID
	}
	if wpsChannelID != "" {
		c.channelType[models.WPS] = wpsChannelID
	}
}

func (c *App) SendMessage(msg, channelID string) {
	if channelID == "" {
		log.Println("Empty channel ID obtained. Could not send message: ", msg)
		return
	}
	if err := c.communicator.SendMessage(msg, channelID); err != nil {
		log.Printf("Could not send message. Error: %s", err)
		if strings.Contains(err.Error(), "message is too long") {
			l := len(msg) / 2
			c.SendMessage(msg[:l], channelID)
			c.SendMessage(msg[l:], channelID)
		} else if strings.Contains(err.Error(), "Too Many Requests") {
			time.Sleep(time.Second * 45)
			c.SendMessage(msg, channelID)
		}
	}
}

func (c *App) Run() error {
	var workerLimit int = 5
	var workerCounter int
	worker := make(chan struct{}, workerLimit)
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM)
	msgs := discord.DiscordChannel

	go c.AutonomousARPScanner()
	go c.AutonomousPortScanner()
	go c.AutonomousWebPortScanner()

	log.Println("Starting application")
	var m models.Message
	for {
		select {
		case m = <-msgs:
			log.Printf("Obtained command - %s from %s", m.Msg, m.ChannelID)
			log.Printf("# of free workers - %d", workerLimit-workerCounter)
			if strings.HasPrefix(m.Msg, "/") {
				if workerCounter < workerLimit {
					workerCounter++
					go func(worker chan struct{}) {
						c.runCommand(m.Msg, m.ChannelID)
						worker <- struct{}{}
					}(worker)
				} else {
					c.SendMessage("I'm too busy already. Try to scan later", m.ChannelID)
				}
			}
		case <-s:
			fmt.Println("\nCtrl+C was pressed. Interrupting the process...")
			close(s)
			close(msgs)
			return nil
		case <-worker:
			workerCounter--
		}
	}
}

func (c *App) runCommand(cmd, channelID string) {
	words := strings.Fields(cmd)
	if len(words) <= 1 {
		if len(words) == 1 {
			c.singleCommandRun(words[0], channelID)
			return
		}
		if err := c.communicator.SendMessage("Not enough arguments", channelID); err != nil {
			log.Printf("Could not send message. Error: %s", err)
		}
		return
	}
	switch words[0] {
	case "/arp_channel_id":
		channelID := words[1]
		c.channelType[models.ARP] = channelID
		msg := fmt.Sprintf("ARP channel ID is set to %s", channelID)
		if err := c.communicator.SendMessage(msg, channelID); err != nil {
			log.Printf("Could not send message. Error %s", err)
		}
	case "/ps_channel_id":
		channelID := words[1]
		c.channelType[models.PS] = channelID
		msg := fmt.Sprintf("PS channel ID is set to %s", channelID)
		if err := c.communicator.SendMessage(msg, channelID); err != nil {
			log.Printf("Could not send message. Error %s", err)
		}
	case "/wps_channel_id":
		channelID := words[1]
		c.channelType[models.WPS] = channelID
		msg := fmt.Sprintf("WPS channel ID is set to %s", channelID)
		if err := c.communicator.SendMessage(msg, channelID); err != nil {
			log.Printf("Could not send message. Error %s", err)
		}
	case "/reply":
		if err := c.communicator.SendMessage(words[1], channelID); err != nil {
			log.Printf("Could not send message. Error %s", err)
		}
	case "/nmap":
		if err := c.AddTargetToNmapScan(words[1], -1); err != nil {
			log.Println("Adding target to nmap scan failed:", err)
		}
	case "/web_nmap":
		if err := c.AddTargetToWebScan(words[1], -1); err != nil {
			log.Println("Adding target to web scan failed:", err)
		}
	case "/arpscan":
		if err := c.AddTargetToARPScan(words[1]); err != nil {
			log.Println("Adding target to arp scan failed:", err)
		}
	}
}

const helpMessage string = `
/get_this_channel_id -> obtains the channel ID where the command is executed
/arp_channel_id -> setting channel ID to send ARP scan results into that channel
/ps_channel_id -> setting channel ID to send all ports scan results into that channel
/wps_channel_id -> setting channel ID to send web ports scan results into that channel
/arpscan -> settings ARP scan target. Accepts both single IP and IP with bitmask
/nmap -> setting all ports scan target. Accepts both single IP and IP with bitmask
/web_nmap -> setting web ports scan target. Accepts both single IP and IP with bitmask
`

func (c *App) singleCommandRun(cmd, channelID string) {
	switch cmd {
	case "/help":
		if err := c.communicator.SendMessage(helpMessage, channelID); err != nil {
			log.Printf("Could not send help message. Error: %s", err)
		}
	case "/get_this_channel_id":
		msg := fmt.Sprintf("This channel ID is %s", channelID)
		if err := c.communicator.SendMessage(msg, channelID); err != nil {
			log.Printf("Could not send message. Error %s", err)
		}
	}
}
