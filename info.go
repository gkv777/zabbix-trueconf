package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// TrueConfInfo ...
type TrueConfInfo struct {
	Online int `json:"online"`
	Busy   int `json:"busy"`
	Active int `json:"active"`
	All    int `json:"all"`
}

// Print ...
func (i *TrueConfInfo) Print() {
	e, err := json.Marshal(i)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(e))
}

// GetUsersInfo ...
func (i *TrueConfInfo) GetUsersInfo(us []User, debug bool) {
	var (
		online = 0
		busy   = 0
		active = 0
	)
	for _, u := range us {
		if u.IsActive == 1 {
			active++
		}
		if u.Status == 2 {
			busy++
			online++
		}
		if u.Status == 1 {
			online++
		}
	}
	if debug {
		log.Printf("GetUsersinfo: %d|%d|%d|%d\n", online, busy, active, len(us))
	}
	i.Active = active
	i.Online = online
	i.Busy = busy
	i.All = len(us)
}
