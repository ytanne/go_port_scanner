package app

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ytanne/go_nessus/pkg/entities"
	m "github.com/ytanne/go_nessus/pkg/models"
)

func (c *App) AddTargetToNmapScan(target string, id int) error {
	log.Println("Obtained all ports scan target", target)
	t, err := c.storage.RetrieveNmapRecord(target, id)

	if err == sql.ErrNoRows {
		log.Printf("No records found for %s", target)
		t, err := c.storage.CreateNewNmapTarget(target, id)
		if err != nil {
			log.Printf("Could not add target %s to the table. Error: %s", target, err)
			return err
		}
		err = c.RunPortScanner(t, -1)
		if err != nil {
			log.Printf("Could not run Nmap scan on %s. Error: %s", t.IP, err)
			t.ErrMsg = err.Error()
			t.ErrStatus = -200
		}
		if _, err := c.storage.SaveNmapResult(t); err != nil {
			log.Println("Storing nmap result failed:", err)
		}

		return nil
	} else if err == nil {
		if time.Since(t.ScanTime) > time.Minute*15 {
			lastResult := len(t.Result)
			err = c.RunPortScanner(t, lastResult)
			if err != nil {
				log.Printf("Could not run Nmap scan on %s. Error: %s", t.IP, err)
				t.ErrMsg = err.Error()
				t.ErrStatus = -200
			}
			if _, err := c.storage.SaveNmapResult(t); err != nil {
				log.Println("Storing nmap result failed:", err)
			}
			return nil
		}

		if t.ErrStatus == -200 {
			c.SendMessage(fmt.Sprintf("Could not do #ALL_PORT scan %s\n%s", t.IP, t.ErrMsg), c.channelType[m.PS])
			return nil
		}
		msg := fmt.Sprintf(
			"%s\nPreviously at #ALL_PORT scan of %s was found:\n%s",
			t.IP,
			t.ScanTime.Format(time.RFC3339),
			t.Result,
		)
		c.SendMessage(msg, c.channelType[m.PS])
		return nil
	}
	log.Printf("Could not retrieve results for %s. Error: %s", target, err)
	return err
}

func (c *App) RunPortScanner(target *entities.NmapTarget, lastResult int) error {
	// c.serv.SendMessage(fmt.Sprintf("Starting #ALL_PORT scanning %s", target.IP))
	ports, err := c.portScanner.ScanPorts(target.IP)
	if err != nil {
		log.Printf("Could not run Port scan on %s. Error: %s", target.IP, err)
		c.SendMessage(fmt.Sprintf("Could not scan #ALL_PORTS of %s", target.IP), c.channelType[m.PS])
		target.ErrMsg = err.Error()
		target.ErrStatus = -200
		return err
	}
	if ports == nil {
		log.Printf("No ports found for %s", target.IP)
		c.SendMessage(fmt.Sprintf("No open #ALL_PORTS of %s found", target.IP), c.channelType[m.PS])
		return nil
	}
	if lastResult != len(ports) {
		if !strings.Contains(target.IP, "/") {
			msg := fmt.Sprintf("Open #ALL_PORTS of %s:\nPORT\tSTATE\tSERVICE\n%s", target.IP, strings.Join(ports, "\n"))
			c.SendMessage(msg, c.channelType[m.PS])
		}
	} else {
		c.SendMessage(fmt.Sprintf("No updates on #ALL_PORTS for %s", target.IP), c.channelType[m.PS])
	}
	target.Result = strings.Join(ports, "; ")
	return nil
}

func (c *App) AutonomousPortScanner() {
	sem := make(chan struct{}, 3)
	ticker := time.NewTicker(time.Minute * 15)
	for ; true; <-ticker.C {
		log.Println("Starting autonomous NMAP check")
		targets, err := c.storage.RetrieveAllNmapTargets()
		if err != nil {
			log.Fatalf("Could not obtain all NMAP targets. Error: %s", err)
		}
		var l int = len(targets)
		log.Printf("There are %d targets for NMAP scan", l)
		for _, target := range targets {
			sem <- struct{}{}
			go func(target *entities.NmapTarget, sem <-chan struct{}) {
				log.Printf("Doing NMAP scan of %s", target.IP)
				lastResult := len(target.Result)
				err = c.RunPortScanner(target, lastResult)
				if err != nil {
					log.Printf("Could not run nmap scan on %s. Error: %s", target.IP, err)
					target.ErrMsg = err.Error()
					target.ErrStatus = -200
				}
				if _, err := c.storage.SaveNmapResult(target); err != nil {
					log.Printf("Could not save ARP result of %s. Error: %s", target.IP, err)
				}
				log.Printf("Finished NMAP scan of %s", target.IP)
				<-sem
			}(target, sem)
		}
		log.Println("Finished autonomous NMAP check. Taking a break")
	}
}
