package app

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ytanne/go_nessus/pkg/entities"
)

var nmapTargets []string

func (c *App) AddTargetToNmapScan(target string, id int) error {
	t, err := c.serv.RetrieveNmapRecord(target, id)

	if err == sql.ErrNoRows {
		log.Printf("No records found for %s", target)
		t, err := c.serv.CreateNewNmapTarget(target, id)
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
		c.serv.SaveNmapResult(t)

		return nil
	} else if err == nil {
		if time.Now().Sub(t.ScanTime) > time.Minute*5 {
			lastResult := len(t.Result)
			err = c.RunPortScanner(t, lastResult)
			if err != nil {
				log.Printf("Could not run Nmap scan on %s. Error: %s", t.IP, err)
				t.ErrMsg = err.Error()
				t.ErrStatus = -200
			}
			if lastResult == -1 || lastResult != len(t.Result) {
				c.serv.SaveNmapResult(t)
			}
			return nil
		}

		if t.ErrStatus == -200 {
			c.SendMessage(fmt.Sprintf("Could not do PORT scan %s\n%s", t.IP, t.ErrMsg))
			return nil
		}

		c.SendMessage(fmt.Sprintf("%s\nPreviously at PORT scan of %s was found:\n%s", t.IP, t.ScanTime.Format(time.RFC3339), t.Result))
		return nil
	}
	log.Printf("Could not retrieve results for %s. Error: %s", target, err)
	return err
}

func (c *App) RunPortScanner(target *entities.NmapTarget, lastResult int) error {
	c.serv.SendMessage(fmt.Sprintf("Starting PORT scanning %s", target.IP))
	ports, err := c.serv.ScanPorts(target.IP)
	if err != nil {
		log.Printf("Error: %s", err)
		c.SendMessage(fmt.Sprintf("Could not scan PORTS of %s", target.IP))
		target.ErrMsg = err.Error()
		target.ErrStatus = -200
		return err
	}
	if ports == "" {
		c.SendMessage(fmt.Sprintf("No open PORTS of %s found", target.IP))
		return nil
	}
	if lastResult == -1 || lastResult != len(target.Result) {
		c.SendMessage(fmt.Sprintf("Open PORTS of %s:\n%s", target.IP, ports))
	} else {
		c.SendMessage(fmt.Sprintf("No updates on PORTS for %s", target.IP))
	}
	target.Result = ports
	return nil
}

func (c *App) AutonomousPortScanner() {
	ticker := time.Tick(time.Minute * 10)
	for {
		targets, err := c.serv.RetrieveOldNmapTargets(10)
		if err != nil {
			log.Printf("Could not retrieve old nmap targets. Error: %s", err)
			continue
		}
		for _, target := range targets {
			lastResult := len(target.Result)
			err = c.RunPortScanner(target, lastResult)
			if err != nil {
				log.Printf("Could not run nmap scan on %s. Error: %s", target.IP, err)
				target.ErrMsg = err.Error()
				target.ErrStatus = -200
			}
			if lastResult == -1 || lastResult != len(target.Result) {
				c.serv.SaveNmapResult(target)
			}
		}
		<-ticker
	}
}
