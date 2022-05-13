package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ytanne/go_nessus/pkg/entities"
	m "github.com/ytanne/go_nessus/pkg/models"
)

func (c *App) AddTargetToWebScan(IP string, id int) error {
	t, err := c.storage.RetrieveWebRecord(context.Background(), IP, id)
	if err == sql.ErrNoRows {
		log.Printf("No records found for %s", IP)

		t, err := c.storage.CreateNewWebTarget(context.Background(), entities.NmapTarget{IP: IP}, id)
		if err != nil {
			log.Printf("Could not add target %s to the table. Error: %s", IP, err)

			return err
		}

		err = c.RunWebPortScanner(t, "")
		if err != nil {
			log.Printf("Could not run Nmap Web scan on %s. Error: %s", t.IP, err)
			t.ErrMsg = err.Error()
			t.ErrStatus = -200
		}

		if _, err := c.storage.SaveWebResult(context.Background(), t); err != nil {
			log.Println("Storing web result failed:", err)
		}

		return nil
	} else if err == nil {
		if time.Since(t.ScanTime) > time.Minute*15 {
			lastResult := t.Result

			err = c.RunWebPortScanner(t, lastResult)
			if err != nil {
				log.Printf("Could not run Nmap Web scan on %s. Error: %s", t.IP, err)
				t.ErrMsg = err.Error()
				t.ErrStatus = -200
			}

			if _, err := c.storage.SaveWebResult(context.Background(), t); err != nil {
				log.Println("Storing web result failed:", err)
			}
			return nil
		}

		if t.ErrStatus == -200 {
			c.SendMessage(fmt.Sprintf("Could not do #WEB_PORT scan %s\n%s", t.IP, t.ErrMsg), c.channelType[m.WPS], startingCount)

			return nil
		}

		msg := fmt.Sprintf(
			"%s\nPreviously at #WEB_PORT scan of %s was found:\n%s",
			t.IP,
			t.ScanTime.Format(time.RFC3339),
			t.Result,
		)

		c.SendMessage(msg, c.channelType[m.WPS], startingCount)
		return nil
	}

	log.Printf("Could not retrieve web results for %s. Error: %s", IP, err)

	return err
}

func (c *App) RunWebPortScanner(target entities.NmapTarget, lastResult string) error {
	ports, err := c.portScanner.ScanWebPorts(c.ctx, target.IP)
	if err != nil {
		log.Printf("Could not run Web Port scan on %s. Error: %s", target.IP, err)
		c.SendMessage(fmt.Sprintf("Could not scan Web_PORTS of %s", target.IP), c.channelType[m.WPS], startingCount)
		target.ErrMsg = err.Error()
		target.ErrStatus = -200

		return err
	}

	if ports == nil {
		log.Printf("No web ports found for %s", target.IP)
		c.SendMessage(fmt.Sprintf("No open #WEB_PORTS of %s found", target.IP), c.channelType[m.WPS], startingCount)

		return nil
	}

	target.Result = strings.Join(ports, "; ")
	if len(lastResult) != len(target.Result) {
		if !strings.Contains(target.IP, "/") {
			msg := fmt.Sprintf(
				"Open #WEB_PORTS of %s:\nPORT\tSTATE\tSERVICE\n%s",
				target.IP,
				ports,
			)

			c.SendMessage(msg, c.channelType[m.WPS], startingCount)
		}
	} else {
		c.SendMessage(fmt.Sprintf("No updates on #WEB_PORTS for %s", target.IP), c.channelType[m.WPS], startingCount)
	}

	return nil
}

func (c *App) AutonomousWebPortScanner() {
	sem := make(chan struct{}, 5)

	ticker := time.NewTicker(time.Minute * 15)
	for ; true; <-ticker.C {
		log.Println("Starting autonomous NMAP Web check")

		targets, err := c.storage.RetrieveAllWebTargets(context.Background())
		if err != nil {
			log.Fatalf("Could not obtain all NMAP web targets. Error: %s", err)
		}

		log.Printf("There are %d targets for NMAP Web scan", len(targets))
		for _, target := range targets {
			sem <- struct{}{}

			go func(target entities.NmapTarget, sem <-chan struct{}) {
				var lastResult string
				log.Printf("Doing NMAP Web scan of %s", target.IP)

				oldTarget, err := c.storage.RetrieveWebRecord(context.Background(), target.IP, target.ARPscanID)
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

				if _, err := c.storage.SaveWebResult(context.Background(), target); err != nil {
					log.Printf("Could not save ARP result of %s. Error: %s", target.IP, err)
				}

				log.Printf("Finished NMAP Web scan of %s", target.IP)
				<-sem
			}(target, sem)
		}

		log.Println("Finished autonomous NMAP web check. Taking a break")
	}
}
