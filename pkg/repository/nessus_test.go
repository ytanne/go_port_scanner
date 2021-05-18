package repository

import (
	"fmt"
	"log"
	"testing"
)

func TestListScans(t *testing.T) {
	nc := NewNessusClient("b6775afa378f5cdce893318e903c1019a8ebd39ea4c4697caf70a3f8e7b97c0a", "773a1ed0a4a9e1329b9ebeeb1929a6742101c9bf49e020a342ff6c41db898005", "https://95.59.127.55:8834/")
	res, err := nc.ListScans()
	if err != nil {
		log.Printf("Could not do list scan: %s", err)
	}
	fmt.Println(res)
}
