package sqlite

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/ytanne/go_nessus/pkg/entities"
)

type DBKeeper interface {
	CreateNewARPTarget(target string) (*entities.ARPTarget, error)
	SaveARPResult(target *entities.ARPTarget) (int, error)
	RetrieveARPRecord(target string) (*entities.ARPTarget, error)
	RetrieveOldARPTargets(timelimit int) ([]*entities.ARPTarget, error)
	RetrieveAllARPTargets() ([]*entities.ARPTarget, error)
	CreateNewNmapTarget(target string, id int) (*entities.NmapTarget, error)
	SaveNmapResult(target *entities.NmapTarget) (int, error)
	RetrieveNmapRecord(target string, id int) (*entities.NmapTarget, error)
	RetrieveOldNmapTargets(timelimit int) ([]*entities.NmapTarget, error)
	RetrieveAllNmapTargets() ([]*entities.NmapTarget, error)
	CreateNewWebTarget(target string, id int) (*entities.NmapTarget, error)
	SaveWebResult(target *entities.NmapTarget) (int, error)
	RetrieveWebRecord(target string, id int) (*entities.NmapTarget, error)
	RetrieveOldWebTargets(timelimit int) ([]*entities.NmapTarget, error)
	RetrieveAllWebTargets() ([]*entities.NmapTarget, error)
}

type database struct {
	db *sql.DB
}

func NewDatabaseRepository(db *sql.DB) DBKeeper {
	return &database{db}
}

func (d *database) CreateNewARPTarget(target string) (*entities.ARPTarget, error) {
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

func (d *database) RetrieveARPRecord(target string) (*entities.ARPTarget, error) {
	var result entities.ARPTarget
	var IPs []byte
	err := d.db.QueryRow(`SELECT id, target, num_of_ips, ips, scan_time, error_status, error_msg FROM arp_targets WHERE target = $1`, target).Scan(&result.ID, &result.Target, &result.NumOfIPs, &IPs, &result.ScanTime, &result.ErrStatus, &result.ErrMsg)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(IPs, &result.IPs); err != nil {
		log.Printf("Could not get IPs for %s", target)
	}
	return &result, nil
}

func (d database) SaveARPResult(target *entities.ARPTarget) (int, error) {
	target.ScanTime = time.Now()
	data, err := json.Marshal(target.IPs)
	if err != nil {
		return -1, err
	}
	_, err = d.db.Exec(`UPDATE arp_targets SET num_of_ips = $1, ips = $2, scan_time = $3, error_status = $4, error_msg = $5 WHERE id = $6 AND target = $7`, target.NumOfIPs, data, target.ScanTime, target.ErrStatus, target.ErrMsg, target.ID, target.Target)
	return target.ID, err
}

func (d database) RetrieveOldARPTargets(timelimit int) ([]*entities.ARPTarget, error) {
	var result []*entities.ARPTarget
	rows, err := d.db.Query(`select id, target, num_of_ips, ips, scan_time, error_status, error_msg from arp_targets where round((julianday(datetime('now')) - julianday(scan_time)) * 1440) > $1`, timelimit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var IPs []byte
	for rows.Next() {
		element := new(entities.ARPTarget)
		rows.Scan(&element.ID, &element.Target, &element.NumOfIPs, &IPs, &element.ScanTime, &element.ErrStatus, &element.ErrMsg)
		if err := json.Unmarshal(IPs, &element.IPs); err != nil {
			log.Printf("Could not unmarshal IPs of %s", element.Target)
		}
		result = append(result, element)
	}
	return result, nil
}

func (d database) RetrieveAllARPTargets() ([]*entities.ARPTarget, error) {
	var result []*entities.ARPTarget
	rows, err := d.db.Query(`select id, target, num_of_ips, ips, scan_time, error_status, error_msg from arp_targets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var IPs []byte
	for rows.Next() {
		element := new(entities.ARPTarget)
		rows.Scan(&element.ID, &element.Target, &element.NumOfIPs, &IPs, &element.ScanTime, &element.ErrStatus, &element.ErrMsg)
		if err := json.Unmarshal(IPs, &element.IPs); err != nil {
			log.Printf("Could not unmarshal IPs of %s", element.Target)
		}
		result = append(result, element)
	}
	return result, nil
}

func (d database) CreateNewNmapTarget(target string, arpId int) (*entities.NmapTarget, error) {
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

func (d database) RetrieveNmapRecord(target string, id int) (*entities.NmapTarget, error) {
	var result entities.NmapTarget
	err := d.db.QueryRow(`SELECT id, arpscan_id, ip, result, scan_time, error_status, error_msg FROM nmap_targets WHERE ip = $1 AND arpscan_id = $2`, target, id).Scan(&result.ID, &result.ARPscanID, &result.IP, &result.Result, &result.ScanTime, &result.ErrStatus, &result.ErrMsg)
	return &result, err
}

func (d database) SaveNmapResult(target *entities.NmapTarget) (int, error) {
	target.ScanTime = time.Now()
	_, err := d.db.Exec(`UPDATE nmap_targets SET result = $1, scan_time = $2, error_status = $3, error_msg = $4 WHERE id = $5 AND ip = $6`, target.Result, target.ScanTime, target.ErrStatus, target.ErrMsg, target.ID, target.IP)
	return target.ID, err
}

func (d database) RetrieveOldNmapTargets(timelimit int) ([]*entities.NmapTarget, error) {
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

func (d database) RetrieveAllNmapTargets() ([]*entities.NmapTarget, error) {
	var result []*entities.NmapTarget
	rows, err := d.db.Query(`select * from nmap_targets`)
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

func (d database) CreateNewWebTarget(target string, arpId int) (*entities.NmapTarget, error) {
	var id int
	_, err := d.db.Exec(`INSERT INTO web_targets (arpscan_id, ip, scan_time) VALUES ($1, $2, $3)`, arpId, target, time.Now())
	if err != nil {
		return nil, err
	}
	err = d.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &entities.NmapTarget{ID: id, ARPscanID: arpId, IP: target}, err
}

func (d database) RetrieveWebRecord(target string, id int) (*entities.NmapTarget, error) {
	var result entities.NmapTarget
	err := d.db.QueryRow(`SELECT id, arpscan_id, ip, result, scan_time, error_status, error_msg FROM web_targets WHERE ip = $1 AND arpscan_id = $2`, target, id).Scan(&result.ID, &result.ARPscanID, &result.IP, &result.Result, &result.ScanTime, &result.ErrStatus, &result.ErrMsg)
	return &result, err
}

func (d database) SaveWebResult(target *entities.NmapTarget) (int, error) {
	target.ScanTime = time.Now()
	_, err := d.db.Exec(`UPDATE web_targets SET result = $1, scan_time = $2, error_status = $3, error_msg = $4 WHERE id = $5 AND ip = $6`, target.Result, target.ScanTime, target.ErrStatus, target.ErrMsg, target.ID, target.IP)
	return target.ID, err
}

func (d database) RetrieveOldWebTargets(timelimit int) ([]*entities.NmapTarget, error) {
	var result []*entities.NmapTarget
	rows, err := d.db.Query(`select * from web_targets where round((julianday(datetime('now')) - julianday(scan_time)) * 1440) > $1 LIMIT 3`, timelimit)
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

func (d database) RetrieveAllWebTargets() ([]*entities.NmapTarget, error) {
	var result []*entities.NmapTarget
	rows, err := d.db.Query(`select * from web_targets`)
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