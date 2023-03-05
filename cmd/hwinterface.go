package main

import (
	"strings"
)

type HWInterface struct {
	clients map[string]*TCPClient
	powerOn bool
	online  bool
}

func NewHWInterface() *HWInterface {
	return &HWInterface{
		clients: map[string]*TCPClient{
			"mpd":  NewTCPClient(TCPClientConf{addr: "192.168.0.95:6600", singleRequest: false}),
			"ctrl": NewTCPClient(TCPClientConf{addr: "192.168.0.95:1025", singleRequest: true}),
		},
	}
}

func (hw *HWInterface) Request(server, command string) ([]string, error) {
	res, err := hw.clients[server].Request(command)
	hw.online = hw.clients["mpd"].online // this one must be always online
	return res, err
}

func (hw *HWInterface) chkPowerState() (bool, error) {
	res, err := hw.Request("ctrl", "check_extpower")
	if err == nil {
		return strings.Split(res[0], ": ")[1] == "1", nil
	}
	return false, err
}
func (hw *HWInterface) togglePower() error {
	_, err := hw.Request("ctrl", "extpower_toggle")
	return err
}
