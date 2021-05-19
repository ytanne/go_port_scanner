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
	t, err := c.serv.RetrieveARPRecord(target)

	if err == sql.ErrNoRows {
		log.Printf("No records found for %s", target)
		t, err := c.serv.CreateNewARPTarget(target)
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
		c.serv.SaveARPResult(t)

		log.Printf("Target ID - %d", t.ID)
		go func() {
			for _, ip := range t.IPs {
				c.serv.CreateNewNmapTarget(ip, t.ID)
			}
		}()

		return nil
	} else if err == nil {
		if time.Now().Sub(t.ScanTime) > time.Minute*5 {
			lastResult := t.IPs
			err = c.RunARPScanner(t, lastResult)
			if err != nil {
				log.Printf("Could not run ARP scan on %s. Error: %s", t.Target, err)
				t.ErrMsg = err.Error()
				t.ErrStatus = -200
			}
			c.serv.SaveARPResult(t)
			return nil
		}

		if t.ErrStatus == -200 {
			c.SendMessage(fmt.Sprintf("Could not ARP scan of %s\n%s", t.Target, t.ErrMsg))
			return nil
		}

		c.SendMessage(fmt.Sprintf("Previously at ARP scan of %s was found %d IPs in %s", t.ScanTime.Format(time.RFC3339), t.NumOfIPs, t.Target))
		return nil
	}
	log.Printf("Could not retrieve records for %s. Error: %s", target, err)
	return err
}

func (c *App) RunARPScanner(target *entities.ARPTarget, lastResult []string) error {
	c.serv.SendMessage(fmt.Sprintf("Starting ARP scanning %s", target.Target))
	ips, err := c.serv.ScanNetwork(target.Target)
	if err != nil {
		c.SendMessage(fmt.Sprintf("Could not do ARP scan network of %s", target.Target))
		return err
	}
	if ips == nil {
		c.SendMessage(fmt.Sprintf("No IPs of %s found in ARP scan", target.Target))
		return nil
	}
	sort.Strings(ips)
	equal := checkIfEqual(lastResult, ips)
	if lastResult == nil || !equal {
		diff := getDifference(lastResult, ips)
		c.SendMessage(fmt.Sprintf("ARP scan: IPs of %s:\n%s", target.Target, strings.Join(diff, "\n")))
	} else {
		c.SendMessage(fmt.Sprintf("No updates for %s on ARP scan", target.Target))
	}
	target.NumOfIPs = len(ips)
	target.IPs = ips
	return nil
}

func (c *App) AutonomousARPScanner() {
	ticker := time.Tick(time.Minute * 10)
	for {
		log.Println("Starting autonomous ARP check")
		targets, err := c.serv.RetrieveOldARPTargets(10)
		if err != nil {
			log.Printf("Could not retrieve old ARP targets. Error: %s", err)
			continue
		}
		log.Printf("Retrieved %d targets for ARP scan", len(targets))
		for _, target := range targets {
			lastResult := target.IPs
			err = c.RunARPScanner(target, lastResult)
			if err != nil {
				log.Printf("Could not run ARP scan on %s. Error: %s", target.Target, err)
				target.ErrMsg = err.Error()
				target.ErrStatus = -200
			}
			c.serv.SaveARPResult(target)
		}
		<-ticker
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
