package main

import (
	"strings"
)

type HWInterface struct {
	clients map[string]*TCPClient
	powerOn bool
}

func NewHWInterface() *HWInterface {
	return &HWInterface{
		clients: map[string]*TCPClient{
			"mpd":  NewClient(TCPClientParms{addr: "192.168.0.95:6600", singleRequest: false}),
			"ctrl": NewClient(TCPClientParms{addr: "192.168.0.95:1025", singleRequest: true}),
		},
	}
}

func (hw *HWInterface) chkPowerStatus() (bool, error) {
	res, err := hw.clients["ctrl"].Request("check_extpower")
	if err == nil {
		return strings.Split(res[0], ": ")[1] == "1", nil
	}
	return false, err
}
func (hw *HWInterface) togglePower() error {
	_, err := hw.clients["ctrl"].Request("extpower_toggle")
	return err
}
