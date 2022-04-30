package models

const (
	_ = iota
	ARP
	PS
	WPS
)

type Message struct {
	Msg       string
	ChannelID string
}
