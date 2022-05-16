package app

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ytanne/go_port_scanner/pkg/models"
)

const (
	startingCount = 0
	sendLimit     = 10
)

type App struct {
	ctx           context.Context
	communicator  Communicator
	storage       Keeper
	portScanner   PortScanner
	nucleiScanner NucleiScanner
	channelType   map[int]string
}

func NewApp(communicator Communicator, storage Keeper, portScanner PortScanner, nucleiScanner NucleiScanner) *App {
	rand.Seed(time.Now().UnixNano())

	return &App{
		communicator:  communicator,
		storage:       storage,
		portScanner:   portScanner,
		nucleiScanner: nucleiScanner,
		channelType:   make(map[int]string),
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

func (c *App) SendMessage(msg, channelID string, counter int) {
	if counter >= sendLimit {
		log.Println("send limit exceeded")

		return
	}

	if channelID == "" {
		log.Println("Empty channel ID obtained. Could not send message: ", msg)

		return
	}

	if err := c.communicator.SendMessage(msg, channelID); err != nil {
		log.Printf("Could not send message. Error: %s", err)

		if strings.Contains(err.Error(), "message is too long") {
			l := len(msg) / 2

			c.SendMessage(msg[:l], channelID, counter+1)
			c.SendMessage(msg[l:], channelID, counter+1)
		} else if strings.Contains(err.Error(), "Too Many Requests") {
			time.Sleep(time.Second * 45)
			c.SendMessage(msg, channelID, counter+1)
		}
	}
}

func (c *App) Run() error {
	var workerLimit int = 5
	var workerCounter int
	worker := make(chan struct{}, workerLimit)
	ctx, cancel := context.WithCancel(context.Background())
	c.ctx = ctx

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM)
	msgChannel := c.communicator.MessageReadChannel()

	go c.AutonomousARPScanner()
	go c.AutonomousPortScanner()
	go c.AutonomousWebPortScanner()

	log.Println("Starting application")
	var m models.Message
	for {
		select {
		case m = <-msgChannel:
			log.Printf("Obtained command - %s from %s", m.Msg, m.ChannelID)
			log.Printf("# of free workers - %d", workerLimit-workerCounter)
			if strings.HasPrefix(m.Msg, "/") {
				if workerCounter < workerLimit {
					workerCounter++
					go func(worker chan struct{}) {
						c.runCommand(m.Msg, m.ChannelID, s)
						worker <- struct{}{}
					}(worker)
				} else {
					c.SendMessage("I'm too busy already. Try to scan later", m.ChannelID, startingCount)
				}
			}
		case <-s:
			fmt.Println("\nCtrl+C was pressed. Interrupting the process...")
			close(s)
			cancel()

			return nil
		case <-worker:
			workerCounter--
		}
	}
}

func (c *App) runCommand(cmd, channelID string, s chan<- os.Signal) {
	words := strings.Fields(cmd)
	if len(words) <= 1 {
		if len(words) == 1 {
			c.singleCommandRun(words[0], channelID, s)
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
	case "/all_nmap":
		msg := fmt.Sprintf("Starting WEB scan of ```%s```", words[1])
		if err := c.communicator.SendMessage(msg, channelID); err != nil {
			log.Printf("Could not send message. Error %s", err)
		}

		if err := c.AddTargetToNmapScan(words[1], -1); err != nil {
			log.Println("Adding target to service scan failed:", err)
		}
	case "/web_nmap":
		msg := fmt.Sprintf("Starting WEB scan of ```%s```", words[1])
		if err := c.communicator.SendMessage(msg, channelID); err != nil {
			log.Printf("Could not send message. Error %s", err)
		}

		if err := c.AddTargetToWebScan(words[1], -1); err != nil {
			log.Println("Adding target to web scan failed:", err)
		}
	case "/arp_scan":
		msg := fmt.Sprintf("Starting ARP scan of ```%s```", words[1])
		if err := c.communicator.SendMessage(msg, channelID); err != nil {
			log.Printf("Could not send message. Error %s", err)
		}

		if err := c.AddTargetToARPScan(words[1]); err != nil {
			log.Println("Adding target to arp scan failed:", err)
		}
	case "/nuclei_scan":
		msg := fmt.Sprintf("Starting nuclei scan of ```%s```", words[1])
		if err := c.communicator.SendMessage(msg, channelID); err != nil {
			log.Printf("Could not send message. Error %s", err)
		}

		var (
			result string
			err    error
		)

		result, err = c.nucleiScanner.ScanURL(context.TODO(), words[1])
		if err != nil {
			log.Println("Adding target to arp scan failed:", err)

			return
		}

		err = c.communicator.SendFile(channelID, result)
		if err != nil {
			log.Printf("Could not send file. Error %s", err)
		}
	}
}

const helpMessage string = `
/help -> print help message
/get_this_channel_id -> obtains the channel ID where the command is executed
/arp_channel_id -> setting channel ID to send ARP scan results into that channel
/ps_channel_id -> setting channel ID to send all ports scan results into that channel
/wps_channel_id -> setting channel ID to send web ports scan results into that channel
/arp_scan -> settings ARP scan target. Accepts both single IP and IP with bitmask
/all_nmap -> setting all ports scan target. Accepts both single IP and IP with bitmask
/web_nmap -> setting web ports scan target. Accepts both single IP and IP with bitmask
/nuclei_scan -> setting nuclei scan target. Accepts either single IP or URL address 
/goodbye -> stop bot
`

func (c *App) singleCommandRun(cmd, channelID string, s chan<- os.Signal) {
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
	case "/goodbye":
		if err := c.communicator.SendMessage("Всем покеда, я спать", channelID); err != nil {
			log.Printf("Could not send message. Error %s", err)
		}
		s <- os.Interrupt
	}
}
