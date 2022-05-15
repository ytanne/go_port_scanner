package app

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ytanne/go_port_scanner/pkg/entities"
	m "github.com/ytanne/go_port_scanner/pkg/models"
)

const (
	nmapScanLimit = 3
)

func (c *App) AddTargetToNmapScan(IP string, id int) error {
	t, err := c.storage.RetrieveNmapRecord(context.Background(), IP, id)
	if t.ID == 0 {
		log.Printf("No records found for %s", IP)

		t, err := c.storage.CreateNewNmapTarget(context.Background(), entities.NmapTarget{IP: IP}, id)
		if err != nil {
			log.Printf("Could not add target %s to the table. Error: %s", IP, err)
			return err
		}

		err = c.RunPortScanner(&t, -1)
		if err != nil {
			log.Printf("Could not run Nmap scan on %s. Error: %s", t.IP, err)
			t.ErrMsg = err.Error()
			t.ErrStatus = -200
		}

		if _, err := c.storage.SaveNmapResult(context.Background(), t); err != nil {
			log.Println("Storing nmap result failed:", err)
		}

		return nil
	} else if err == nil {
		if time.Since(t.ScanTime) > time.Minute*15 {
			lastResult := len(t.Result)

			err = c.RunPortScanner(&t, lastResult)
			if err != nil {
				log.Printf("Could not run Nmap scan on %s. Error: %s", t.IP, err)
				t.ErrMsg = err.Error()
				t.ErrStatus = -200
			}

			if _, err := c.storage.SaveNmapResult(context.Background(), t); err != nil {
				log.Println("Storing nmap result failed:", err)
			}

			return nil
		}

		if t.ErrStatus == -200 {
			c.SendMessage(fmt.Sprintf("Could not do #ALL_PORT scan %s\n%s", t.IP, t.ErrMsg), c.channelType[m.PS], startingCount)

			return nil
		}

		msg := fmt.Sprintf(
			"%s\nPreviously at #ALL_PORT scan of %s was found:\n%s",
			t.IP,
			t.ScanTime.Format(time.RFC3339),
			t.Result,
		)

		c.SendMessage(msg, c.channelType[m.PS], startingCount)
		return nil
	}

	log.Printf("Could not retrieve results for %s. Error: %s", IP, err)
	return err
}

func (c *App) RunPortScanner(target *entities.NmapTarget, lastResult int) error {
	ports, err := c.portScanner.ScanPorts(c.ctx, target.IP)
	if err != nil {
		log.Printf("Could not run Port scan on %s. Error: %s", target.IP, err)
		c.SendMessage(fmt.Sprintf("Could not scan #ALL_PORTS of %s", target.IP), c.channelType[m.PS], startingCount)
		target.ErrMsg = err.Error()
		target.ErrStatus = -200

		return err
	}

	if ports == nil {
		log.Printf("No ports found for %s", target.IP)
		c.SendMessage(fmt.Sprintf("No open #ALL_PORTS of %s found", target.IP), c.channelType[m.PS], startingCount)

		return nil
	}

	if lastResult != len(ports) {
		if !strings.Contains(target.IP, "/") {
			msg := fmt.Sprintf("Open #ALL_PORTS of %s:\nPORT\tSTATE\tSERVICE\n%s", target.IP, strings.Join(ports, "\n"))
			c.SendMessage(msg, c.channelType[m.PS], startingCount)
		}
	} else {
		c.SendMessage(fmt.Sprintf("No updates on #ALL_PORTS for %s", target.IP), c.channelType[m.PS], startingCount)
	}

	target.Result = strings.Join(ports, "; ")
	return nil
}

func (c *App) AutonomousPortScanner() {
	sem := make(chan struct{}, nmapScanLimit)
	scanInterval := time.Minute * 15

	ticker := time.NewTicker(scanInterval)
	for ; true; <-ticker.C {
		log.Println("Starting autonomous NMAP check")
		targets, err := c.storage.RetrieveOldNmapTargets(context.Background(), int(scanInterval.Minutes()))
		if err != nil {
			log.Printf("could not obtain old NMAP targets. Error: %s", err)

			continue
		}

		l := len(targets)
		log.Printf("There are %d targets for NMAP scan", l)
		for _, target := range targets {
			sem <- struct{}{}

			go func(target entities.NmapTarget, sem <-chan struct{}) {
				log.Printf("Doing NMAP scan of %s", target.IP)
				lastResult := len(target.Result)

				err = c.RunPortScanner(&target, lastResult)
				if err != nil {
					log.Printf("Could not run nmap scan on %s. Error: %s", target.IP, err)
					target.ErrMsg = err.Error()
					target.ErrStatus = -200
				}

				if _, err := c.storage.SaveNmapResult(context.Background(), target); err != nil {
					log.Printf("Could not save ARP result of %s. Error: %s", target.IP, err)
				}

				log.Printf("Finished NMAP scan of %s", target.IP)
				<-sem
			}(target, sem)
		}

		log.Println("Finished autonomous NMAP check. Taking a break")
	}
}
