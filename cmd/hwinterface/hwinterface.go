package hwinterface

import (
	"remotecc/cmd/tcpclient"
	"strings"
)

type HWInterface struct {
	clients map[string]*tcpclient.Client
	PowerOn bool
	Online  bool
}

func NewHWInterface() *HWInterface {
	return &HWInterface{
		clients: map[string]*tcpclient.Client{
			"mpd":  tcpclient.New(tcpclient.Conf{Addr: "192.168.0.95:6600", SingleRequest: false}),
			"ctrl": tcpclient.New(tcpclient.Conf{Addr: "192.168.0.95:1025", SingleRequest: true}),
		},
	}
}

func (hw *HWInterface) Request(server, command string) ([]string, error) {
	res, err := hw.clients[server].Request(command)
	hw.Online = hw.clients["mpd"].Online // this one must be always online
	return res, err
}

func (hw *HWInterface) ChkPowerState() (bool, error) {
	res, err := hw.Request("ctrl", "check_extpower")
	if err == nil {
		return strings.Split(res[0], ": ")[1] == "1", nil
	}
	return false, err
}
func (hw *HWInterface) TogglePower() error {
	_, err := hw.Request("ctrl", "extpower_toggle")
	return err
}