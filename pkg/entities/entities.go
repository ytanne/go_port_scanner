package entities

import "time"

type ARPTarget struct {
	ID        int `bson:"id"`
	Target    string `bson:"target"`
	NumOfIPs  int `bson:"num_of_ips"`
	IPs       []string `bson:"ips"`
	ScanTime  time.Time  `bson:"scan_time"`
	ErrStatus int `bson:"err_status"`
	ErrMsg    string `bson:"err_msg"`
}

type NmapTarget struct {
	ID        int `bson:"id"`
	ARPscanID int `bson:"arp_scan_id"`
	IP        string `bson:"ip"`
	Result    string `bson:"result"`
	ScanTime  time.Time `bson:"scan_time"`
	ErrStatus int `bson:"err_status"`
	ErrMsg    string `bson:"err_msg"`
}

type Resource struct {
	ID   int    `json:"int"`
	Name string `json:"string"`
	Type string `json:"trash"`
}

type ScanDetails struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Enabled string `json:"enabled"`
}

type ScanList struct {
	Folders []Resource    `json:"folders"`
	Scans   []ScanDetails `json:"resource"`
}
