package repository

import (
	"database/sql"
	"log"
	"time"

	"github.com/ytanne/go_nessus/pkg/entities"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{db}
}

func (d *Database) CreateNewARPTarget(target string) (*entities.ARPTarget, error) {
	res, err := d.db.Exec(`INSERT INTO arp_targets (target) VALUES ($1)`, target)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	log.Printf("Last ID - %d", id)
	return &entities.ARPTarget{ID: int(id), Target: target}, err
}

func (d *Database) RetrieveARPRecord(target string) (*entities.ARPTarget, error) {
	var result entities.ARPTarget
	err := d.db.QueryRow(`SELECT id, target, num_of_ips, scan_time, error_status, error_msg FROM arp_targets WHERE target = $1`, target).Scan(&result.ID, &result.Target, &result.NumOfIPs, &result.ScanTime, &result.ErrStatus, &result.ErrMsg)
	return &result, err
}

func (d *Database) SaveARPResult(target *entities.ARPTarget) (int, error) {
	target.ScanTime = time.Now()
	_, err := d.db.Exec(`UPDATE arp_targets SET num_of_ips = $1, scan_time = $2, error_status = $3, error_msg = $4 WHERE id = $5 AND target = $6`, target.NumOfIPs, target.ScanTime, target.ErrStatus, target.ErrMsg, target.ID, target.Target)
	return target.ID, err
}

func (d *Database) RetrieveOldARPTargets(timelimit int) ([]*entities.ARPTarget, error) {
	var result []*entities.ARPTarget
	rows, err := d.db.Query(`select * from arp_targets where round((julianday(datetime('now')) - julianday(scan_time)) * 1440) > $1`, timelimit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		element := new(entities.ARPTarget)
		rows.Scan(&element.ID, &element.Target, &element.NumOfIPs, &element.ScanTime, &element.ErrStatus, &element.ErrMsg)
		result = append(result, element)
	}
	return result, nil
}

func (d *Database) CreateNewNmapTarget(target string, arpId int) (*entities.NmapTarget, error) {
	var id int
	_, err := d.db.Exec(`INSERT INTO nmap_targets (arpscan_id, ip, scan_time) VALUES ($1, $2, $3)`, arpId, target, time.Now())
	if err != nil {
		return nil, err
	}
	err = d.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &entities.NmapTarget{ID: id, ARPscanID: arpId, IP: target}, err
}

func (d *Database) RetrieveNmapRecord(target string, id int) (*entities.NmapTarget, error) {
	var result entities.NmapTarget
	err := d.db.QueryRow(`SELECT id, arpscan_id, ip, result, scan_time, error_status, error_msg FROM nmap_targets WHERE ip = $1 AND arpscan_id = $2`, target, id).Scan(&result.ID, &result.ARPscanID, &result.IP, &result.Result, &result.ScanTime, &result.ErrStatus, &result.ErrMsg)
	return &result, err
}

func (d *Database) SaveNmapResult(target *entities.NmapTarget) (int, error) {
	target.ScanTime = time.Now()
	_, err := d.db.Exec(`UPDATE nmap_targets SET result = $1, scan_time = $2, error_status = $3, error_msg = $4 WHERE id = $5 AND ip = $6`, target.Result, target.ScanTime, target.ErrStatus, target.ErrMsg, target.ID, target.IP)
	return target.ID, err
}

func (d *Database) RetrieveOldNmapTargets(timelimit int) ([]*entities.NmapTarget, error) {
	var result []*entities.NmapTarget
	rows, err := d.db.Query(`select * from nmap_targets where round((julianday(datetime('now')) - julianday(scan_time)) * 1440) > $1 LIMIT 3`, timelimit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		element := new(entities.NmapTarget)
		rows.Scan(&element.ID, &element.ARPscanID, &element.IP, &element.Result, &element.ScanTime, &element.ErrStatus, &element.ErrMsg)
		result = append(result, element)
	}
	return result, nil
}
