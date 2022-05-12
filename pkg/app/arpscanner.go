package app

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/ytanne/go_nessus/pkg/entities"
	m "github.com/ytanne/go_nessus/pkg/models"
)

func (c *App) AddTargetToARPScan(target string) error {
	t, err := c.storage.RetrieveARPRecord(target)
	if err == sql.ErrNoRows {
		log.Printf("No records found for %s", target)
		t, err := c.storage.CreateNewARPTarget(target)
		if err != nil {
			log.Printf("Could not add target %s to the table. Error: %s", target, err)

			return err
		}

		err = c.RunARPScanner(t, nil)
		if err != nil {
			log.Printf("Could not run ARP scan on %s. Error: %s", t.Target, err)
			t.ErrMsg = err.Error()
			t.ErrStatus = -200
		}

		if _, err := c.storage.SaveARPResult(t); err != nil {
			log.Println("Storing ARP result failed:", err)

			return err
		}

		log.Printf("Target ID - %d", t.ID)
		go func() {
			for _, ip := range t.IPs {
				_, err := c.storage.CreateNewNmapTarget(ip, t.ID)
				if err != nil {
					log.Println("Creating new nmap target failed")
				}

				_, err = c.storage.CreateNewWebTarget(ip, t.ID)
				if err != nil {
					log.Println("Creating new web target failed")
				}
			}
		}()

		return nil
	} else if err == nil {
		if time.Since(t.ScanTime) > time.Minute*5 {
			lastResult := t.IPs

			err = c.RunARPScanner(t, lastResult)
			if err != nil {
				log.Printf("Could not run ARP scan on %s. Error: %s", t.Target, err)
				t.ErrMsg = err.Error()
				t.ErrStatus = -200
			}

			if _, err := c.storage.SaveARPResult(t); err != nil {
				log.Println("Storing ARP result failed:", err)
			}

			return nil
		}

		if t.ErrStatus == -200 {
			c.SendMessage(
				fmt.Sprintf("Could not ARP scan of %s\n%s", t.Target, t.ErrMsg),
				c.channelType[m.ARP],
				startingCount,
			)
			return nil
		}
		msg := fmt.Sprintf(
			"Previously at ARP scan of %s was found %d IPs in %s\n%s",
			t.ScanTime.Format(time.RFC3339),
			t.NumOfIPs,
			t.Target,
			strings.Join(t.IPs, "\n"),
		)

		c.SendMessage(msg, c.channelType[m.ARP], startingCount)
		return nil
	}

	log.Printf("Could not retrieve records for %s. Error: %s", target, err)
	return err
}

func (c *App) RunARPScanner(target *entities.ARPTarget, lastResult []string) error {
	ips, err := c.portScanner.ScanNetwork(c.ctx, target.Target)
	if err != nil {
		c.SendMessage(
			fmt.Sprintf("Could not do ARP scan network of %s", target.Target),
			c.channelType[m.ARP],
			startingCount,
		)
		return err
	}

	if ips == nil {
		c.SendMessage(
			fmt.Sprintf("No IPs of %s found in ARP scan", target.Target),
			c.channelType[m.ARP],
			startingCount,
		)
		return nil
	}

	sort.Strings(ips)
	if lastResult == nil {
		msg := fmt.Sprintf(
			"ARP scan of %s:\n%s",
			target.Target,
			strings.Join(ips, "\n"),
		)

		c.SendMessage(msg, c.channelType[m.ARP], startingCount)
	}

	equal := checkIfEqual(lastResult, ips)
	if !equal {
		diff := getDifference(lastResult, ips)
		if len(lastResult) > len(ips) {
			msg := fmt.Sprintf(
				"ARP scan of %s. No more available:\n%s",
				target.Target,
				strings.Join(diff, "\n"),
			)
			c.SendMessage(msg, c.channelType[m.ARP], startingCount)
		} else {
			msg := fmt.Sprintf(
				"ARP scan of %s. New IPs detected:\n%s",
				target.Target,
				strings.Join(diff, "\n"),
			)
			c.SendMessage(msg, c.channelType[m.ARP], startingCount)
		}
	} else {
		c.SendMessage(
			fmt.Sprintf("No updates for %s on ARP scan", target.Target),
			c.channelType[m.ARP],
			startingCount,
		)
	}

	target.NumOfIPs = len(ips)
	target.IPs = ips

	return nil
}

func (c *App) AutonomousARPScanner() {
	sem := make(chan struct{}, 2)

	ticker := time.NewTicker(time.Minute * 15)
	for ; true; <-ticker.C {
		log.Println("Starting autonomous ARP check")

		targets, err := c.storage.RetrieveAllARPTargets()
		if err != nil {
			log.Fatalf("Could not obtain all ARP targets: %s", err)
		}

		log.Printf("There are %d targets for ARP scan", len(targets))

		for _, target := range targets {
			sem <- struct{}{}

			go func(target *entities.ARPTarget, sem <-chan struct{}) {
				lastResult := target.IPs

				err = c.RunARPScanner(target, lastResult)
				if err != nil {
					log.Printf("Could not run ARP scan on %s. Error: %s", target.Target, err)
					target.ErrMsg = err.Error()
					target.ErrStatus = -200
				}

				if _, err := c.storage.SaveARPResult(target); err != nil {
					log.Printf("Could not save ARP result of %s. Error: %s", target.Target, err)
				}

				<-sem
			}(target, sem)
		}

		log.Println("Finished autonomous ARP check. Taking a break")
	}
}

func checkIfEqual(arr1 []string, arr2 []string) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	sort.Strings(arr1)
	sort.Strings(arr2)
	for i := 0; i < len(arr1); i++ {
		if arr1[i] != arr2[i] {
			return false
		}
	}

	return true
}

func getDifference(slice1 []string, slice2 []string) []string {
	original := make(map[string]int)
	for _, ip := range slice1 {
		original[ip] = 1
	}

	for _, ip := range slice2 {
		original[ip] = 2
	}

	var diff []string
	for ip, val := range original {
		if val == 1 {
			diff = append(diff, ip)
		}
	}

	return diff
}
