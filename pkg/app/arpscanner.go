package app

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/ytanne/go_nessus/pkg/entities"
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
			c.SendMessage(fmt.Sprintf("Could not ARP scan of %s\n%s", t.Target, t.ErrMsg), c.arpChannelID)
			return nil
		}
		msg := fmt.Sprintf(
			"Previously at ARP scan of %s was found %d IPs in %s\n%s",
			t.ScanTime.Format(time.RFC3339),
			t.NumOfIPs,
			t.Target,
			strings.Join(t.IPs, "\n"),
		)
		c.SendMessage(msg, c.arpChannelID)
		return nil
	}
	log.Printf("Could not retrieve records for %s. Error: %s", target, err)
	return err
}

func (c *App) RunARPScanner(target *entities.ARPTarget, lastResult []string) error {
	// c.serv.SendMessage(fmt.Sprintf("Starting ARP scanning %s", target.Target))
	ips, err := c.portScanner.ScanNetwork(target.Target)
	if err != nil {
		c.SendMessage(fmt.Sprintf("Could not do ARP scan network of %s", target.Target), c.arpChannelID)
		return err
	}
	if ips == nil {
		c.SendMessage(fmt.Sprintf("No IPs of %s found in ARP scan", target.Target), c.arpChannelID)
		return nil
	}
	sort.Strings(ips)
	equal := checkIfEqual(lastResult, ips)
	if lastResult == nil || !equal {
		diff := getDifference(lastResult, ips)
		if lastResult == nil {
			log.Printf("Last result for %s is nil", target.Target)
		} else if len(lastResult) > len(ips) {
			msg := fmt.Sprintf(
				"ARP scan of %s:\n%s\nNo more available:\n%s",
				target.Target,
				strings.Join(ips, "\n"),
				strings.Join(diff, "\n"),
			)
			c.SendMessage(msg, c.arpChannelID)
		} else {
			msg := fmt.Sprintf(
				"ARP scan of %s:\n%s\nNew IPs detected:\n%s",
				target.Target,
				strings.Join(lastResult, "\n"),
				strings.Join(diff, "\n"),
			)
			c.SendMessage(msg, c.arpChannelID)
		}
	} else {
		c.SendMessage(fmt.Sprintf("No updates for %s on ARP scan", target.Target), c.arpChannelID)
	}
	target.NumOfIPs = len(ips)
	target.IPs = ips
	return nil
}

func (c *App) AutonomousARPScanner() {
	targets, err := c.storage.RetrieveAllARPTargets()
	if err != nil {
		log.Fatalf("Could not obtain all ARP targets: %s", err)
	}
	sem := make(chan struct{}, 2)
	ticker := time.Tick(time.Minute * 15)
	for range ticker {
		log.Println("Starting autonomous ARP check")
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
	var diff []string

	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			if !found {
				diff = append(diff, s1)
			}
		}
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

// func (c App) goThroughTargets() {
// 	log.Printf("Going through targets. Obtained length: %d", len(arpTargets))
// 	for _, target := range arpTargets {
// 		if err := c.RunARPScanner(target); err != nil {
// 			log.Printf("Could not scan %s in autonomous mode: %s", target, err)
// 		}
// 	}
// }
