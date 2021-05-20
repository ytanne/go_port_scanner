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

func (c *App) AddTargetToWebScan(target string, id int) error {
	t, err := c.serv.RetrieveWebRecord(target, id)

	if err == sql.ErrNoRows {
		log.Printf("No records found for %s", target)
		t, err := c.serv.CreateNewWebTarget(target, id)
		if err != nil {
			log.Printf("Could not add target %s to the table. Error: %s", target, err)
			return err
		}
		err = c.RunWebPortScanner(t, "")
		if err != nil {
			log.Printf("Could not run Nmap Web scan on %s. Error: %s", t.IP, err)
			t.ErrMsg = err.Error()
			t.ErrStatus = -200
		}
		c.serv.SaveWebResult(t)

		return nil
	} else if err == nil {
		if time.Now().Sub(t.ScanTime) > time.Minute*5 {
			lastResult := t.Result
			err = c.RunWebPortScanner(t, lastResult)
			if err != nil {
				log.Printf("Could not run Nmap Web scan on %s. Error: %s", t.IP, err)
				t.ErrMsg = err.Error()
				t.ErrStatus = -200
			}
			c.serv.SaveWebResult(t)
			return nil
		}

		if t.ErrStatus == -200 {
			c.SendMessage(fmt.Sprintf("Could not do #WEB_PORT scan %s\n%s", t.IP, t.ErrMsg))
			return nil
		}

		c.SendMessage(fmt.Sprintf("%s\nPreviously at #WEB_PORT scan of %s was found:\n%s", t.IP, t.ScanTime.Format(time.RFC3339), t.Result))
		return nil
	}
	log.Printf("Could not retrieve web results for %s. Error: %s", target, err)
	return err
}

func (c *App) RunWebPortScanner(target *entities.NmapTarget, lastResult string) error {
	// c.serv.SendMessage(fmt.Sprintf("Starting #WEB_PORT scanning %s", target.IP))
	ports, err := c.serv.ScanWebPorts(target.IP)
	if err != nil {
		log.Printf("Could not run Web Port scan on %s. Error: %s", target.IP, err)
		c.SendMessage(fmt.Sprintf("Could not scan Web_PORTS of %s", target.IP))
		target.ErrMsg = err.Error()
		target.ErrStatus = -200
		return err
	}
	if ports == nil {
		log.Printf("No web ports found for %s", target.IP)
		c.SendMessage(fmt.Sprintf("No open #WEB_PORTS of %s found", target.IP))
		return nil
	}
	target.Result = strings.Join(ports, "; ")
	if lastResult != target.Result {
		c.SendMessage(fmt.Sprintf("Open #WEB_PORTS of %s:\nPORT\tSTATE\tSERVICE\n%s", target.IP, strings.Join(ports, "\n")))
	} else {
		c.SendMessage(fmt.Sprintf("No updates on #WEB_PORTS for %s", target.IP))
	}
	return nil
}

func (c *App) AutonomousWebPortScanner() {
	targets, err := c.serv.RetrieveAllWebTargets()
	if err != nil {
		log.Fatalf("Could not obtain all NMAP web targets. Error: %s", err)
	}
	var wg sync.WaitGroup
	var l int = len(targets)

	for {
		log.Println("Starting autonomous NMAP Web check")
		log.Printf("There are %d targets for NMAP Web scan", l)
		for i, target := range targets {
			wg.Add(1)
			go func(target *entities.NmapTarget) {
				var lastResult string
				log.Printf("Doing NMAP Web scan of %s", target.IP)

				oldTarget, err := c.serv.RetrieveWebRecord(target.IP, target.ID)
				if err != nil {
					log.Printf("Could not obtain old web record. Error: %s", err)
					log.Printf("IP: %s. ID: %d", target.IP, target.ID)
				} else {
					lastResult = oldTarget.Result
				}
				err = c.RunWebPortScanner(target, lastResult)
				if err != nil {
					log.Printf("Could not run nmap web scan on %s. Error: %s", target.IP, err)
					target.ErrMsg = err.Error()
					target.ErrStatus = -200
				}
				if _, err := c.serv.SaveWebResult(target); err != nil {
					log.Printf("Could not save ARP result of %s. Error: %s", target.IP, err)
				}
				log.Printf("Finished NMAP Web scan of %s", target.IP)
				wg.Done()
			}(target)
			if (i+1)%5 == 0 || (i+1) == l {
				wg.Wait()
			}
		}
		log.Println("Finished autonomous NMAP web check. Taking a break")
		time.Sleep(time.Minute * 5)
	}
}
