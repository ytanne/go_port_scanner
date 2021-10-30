package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ytanne/go_nessus/pkg/service/nmap"
	"github.com/ytanne/go_nessus/pkg/service/sqlite"
	"github.com/ytanne/go_nessus/pkg/service/telegram"
)

type App struct {
	communicator telegram.Communicator
	storage      sqlite.DBKeeper
	portScanner  nmap.NmapScanner
}

func NewApp(communicator telegram.Communicator, storage sqlite.DBKeeper, portScanner nmap.NmapScanner) *App {
	return &App{communicator, storage, portScanner}
}

func (c *App) SendMessage(msg string) {
	if err := c.communicator.SendMessage(msg); err != nil {
		log.Printf("Could not send message. Error: %s", err)
		if strings.Contains(err.Error(), "message is too long") {
			l := len(msg) / 2
			c.SendMessage(msg[:l])
			c.SendMessage(msg[l:])
		} else if strings.Contains(err.Error(), "Too Many Requests") {
			time.Sleep(time.Second * 45)
			c.SendMessage(msg)
		}
	}
}

func (c *App) Run() error {
	var workerLimit int = 3
	var workerCounter int
	worker := make(chan struct{}, workerLimit)
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM)
	msgs := make(chan string, 3)

	go func() {
		err := c.communicator.ReadMessage(msgs)
		if err != nil {
			log.Printf("Could not read message from the bot. Error: %s", err)
		}
	}()

	go c.AutonomousARPScanner()
	go c.AutonomousPortScanner()
	go c.AutonomousWebPortScanner()

	var cmd string
	for {
		select {
		case cmd = <-msgs:
			log.Printf("Obtained command - %s", cmd)
			log.Printf("# of free workers - %d", workerLimit-workerCounter)
			if strings.HasPrefix(cmd, "/") {
				if workerCounter < workerLimit {
					workerCounter++
					go func(worker chan struct{}) {
						c.runCommand(cmd)
						worker <- struct{}{}
					}(worker)
				} else {
					c.SendMessage("I'm too busy already. Try to scan later")
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

func (c *App) runCommand(cmd string) {
	words := strings.Fields(cmd)
	if len(words) <= 1 {
		if err := c.communicator.SendMessage("Not enough arguments"); err != nil {
			log.Printf("Could not send message. Error: %s", err)
		}
		return
	}
	switch words[0] {
	case "/reply":
		if err := c.communicator.SendMessage(words[1]); err != nil {
			log.Printf("Could not send message. Error %s", err)
		}
	case "/nmap":
		c.AddTargetToNmapScan(words[1], -1)
	case "/web_nmap":
		c.AddTargetToWebScan(words[1], -1)
	case "/arpscan":
		c.AddTargetToARPScan(words[1])
	}
}
