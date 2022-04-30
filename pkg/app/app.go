package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ytanne/go_nessus/pkg/service"
)

const (
	sendLimit     = 5
	startingCount = 0
)

type App struct {
	ctx  context.Context
	serv *service.Service
}

func NewApp(serv *service.Service) *App {
	return &App{
		serv: serv,
	}
}

func (c *App) SendMessage(msg string, counter int) {
	if counter >= sendLimit {
		log.Println("send limit exceeded")

		return
	}

	if err := c.serv.SendMessage(msg); err != nil {
		if strings.Contains(err.Error(), "message is too long") {
			l := len(msg) / 2
			time.Sleep(time.Second * 5)
			c.SendMessage(msg[:l], counter+1)
			c.SendMessage(msg[l:], counter+1)
		} else if strings.Contains(err.Error(), "Too Many Requests") {
			time.Sleep(time.Second * 45)
			c.SendMessage(msg, counter+1)
		}
	}
}

func (c *App) Run() error {
	var workerLimit int = 3
	var workerCounter int
	worker := make(chan struct{}, workerLimit)
	ctx, cancel := context.WithCancel(context.Background())

	c.ctx = ctx

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM)
	msgs := make(chan string, 3)

	go func() {
		err := c.serv.ReadMessage(ctx, msgs)
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
					c.SendMessage("I'm too busy already. Try to scan later", startingCount)
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

func (c *App) runCommand(cmd string) {
	words := strings.Fields(cmd)
	if len(words) <= 1 {
		if err := c.serv.SendMessage("Not enough arguments"); err != nil {
			log.Printf("Could not send message. Error: %s", err)
		}

		return
	}

	switch words[0] {
	case "/reply":
		if err := c.serv.SendMessage(words[1]); err != nil {
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
