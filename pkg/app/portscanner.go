package app

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
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
			c.serv.SaveNmapResult(t)
			return nil
		}

		if t.ErrStatus == -200 {
			c.SendMessage(fmt.Sprintf("Could not do#ALL_PORT scan %s\n%s", t.IP, t.ErrMsg))
			return nil
		}

		c.SendMessage(fmt.Sprintf("%s\nPreviously at#ALL_PORT scan of %s was found:\n%s", t.IP, t.ScanTime.Format(time.RFC3339), t.Result))
		return nil
	}
	log.Printf("Could not retrieve results for %s. Error: %s", target, err)
	return err
}

func (c *App) RunPortScanner(target *entities.NmapTarget, lastResult int) error {
	c.serv.SendMessage(fmt.Sprintf("Starting #ALL_PORT scanning %s", target.IP))
	ports, err := c.serv.ScanPorts(target.IP)
	if err != nil {
		log.Printf("Could not run Port scan on %s. Error: %s", target.IP, err)
		c.SendMessage(fmt.Sprintf("Could not scan #ALL_PORTS of %s", target.IP))
		target.ErrMsg = err.Error()
		target.ErrStatus = -200
		return err
	}
	if ports == nil {
		log.Printf("No ports found for %s", target.IP)
		c.SendMessage(fmt.Sprintf("No open #ALL_PORTS of %s found", target.IP))
		return nil
	}
	if lastResult != len(ports) {
		c.SendMessage(fmt.Sprintf("Open #ALL_PORTS of %s:\nPORT\tSTATE\tSERVICE\n%s", target.IP, strings.Join(ports, "\n")))
	} else {
		c.SendMessage(fmt.Sprintf("No updates on #ALL_PORTS for %s", target.IP))
	}
	target.Result = strings.Join(ports, "; ")
	return nil
}

func (c *App) RunWebPortScanner(target *entities.NmapTarget, lastResult int) error {
	c.serv.SendMessage(fmt.Sprintf("Starting #WEB_PORT scanning %s", target.IP))
	ports, err := c.serv.ScanWebPorts(target.IP)
	if err != nil {
		log.Printf("Could not run Web Port scan on %s. Error: %s", target.IP, err)
		c.SendMessage(fmt.Sprintf("Could not scan Web_PORTS of %s", target.IP))
		target.ErrMsg = err.Error()
		target.ErrStatus = -200
		return err
	}
	if ports == nil {
		log.Printf("No ports found for %s", target.IP)
		c.SendMessage(fmt.Sprintf("No open #WEB_PORTS of %s found", target.IP))
		return nil
	}
	if lastResult != len(ports) {
		c.SendMessage(fmt.Sprintf("Open #WEB_PORTS of %s:\nPORT\tSTATE\tSERVICE\n%s", target.IP, strings.Join(ports, "\n")))
	} else {
		c.SendMessage(fmt.Sprintf("No updates on #WEB_PORTS for %s", target.IP))
	}
	target.Result = strings.Join(ports, "; ")
	return nil
}

func (c *App) AutonomousPortScanner() {
	targets, err := c.serv.RetrieveAllNmapTargets()
	if err != nil {
		log.Fatalf("Could not obtain all NMAP targets. Error: %s", err)
	}
	var wg sync.WaitGroup
	var l int = len(targets)
	for {
		log.Println("Starting autonomous NMAP check")
		log.Printf("There are %d targets for NMAP scan", l)
		for i, target := range targets {
			wg.Add(1)
			go func(target *entities.NmapTarget) {
				log.Printf("Doing NMAP scan of %s", target.IP)
				lastResult := len(target.Result)
				err = c.RunPortScanner(target, lastResult)
				if err != nil {
					log.Printf("Could not run nmap scan on %s. Error: %s", target.IP, err)
					target.ErrMsg = err.Error()
					target.ErrStatus = -200
				}
				c.serv.SaveNmapResult(target)
				log.Printf("Finished NMAP scan of %s", target.IP)
				wg.Done()
			}(target)
			if (i+1)%3 == 0 || (i+1) == l {
				wg.Wait()
			}
		}
		log.Println("Finished autonomous ARP check. Taking a break")
		time.Sleep(time.Minute * 5)
	}
}

func (c *App) AutonomousWebPortScanner() {
	targets, err := c.serv.RetrieveAllNmapTargets()
	if err != nil {
		log.Fatalf("Could not obtain all NMAP targets. Error: %s", err)
	}
	var wg sync.WaitGroup
	var l int = len(targets)
	for {
		log.Println("Starting autonomous NMAP Web check")
		log.Printf("There are %d targets for NMAP Web scan", l)
		for i, target := range targets {
			wg.Add(1)
			go func(target *entities.NmapTarget) {
				log.Printf("Doing NMAP Web scan of %s", target.IP)
				lastResult := len(target.Result)
				err = c.RunWebPortScanner(target, lastResult)
				if err != nil {
					log.Printf("Could not run nmap scan on %s. Error: %s", target.IP, err)
					target.ErrMsg = err.Error()
					target.ErrStatus = -200
				}
				log.Printf("Finished NMAP scan of %s", target.IP)
				wg.Done()
			}(target)
			if (i+1)%5 == 0 || (i+1) == l {
				wg.Wait()
			}
		}
		log.Println("Finished autonomous ARP check. Taking a break")
		time.Sleep(time.Minute * 5)
	}
}
