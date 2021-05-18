package entities

import "time"

type ARPTarget struct {
	ID        int
	Target    string
	NumOfIPs  int
	IPs       []string
	ScanTime  time.Time
	ErrStatus int
	ErrMsg    string
}

type NmapTarget struct {
	ID        int
	ARPscanID int
	IP        string
	Result    string
	ScanTime  time.Time
	ErrStatus int
	ErrMsg    string
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
