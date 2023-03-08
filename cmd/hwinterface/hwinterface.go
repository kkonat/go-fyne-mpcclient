package hwinterface

import (
	"fmt"
	"remotecc/cmd/storage"
	"remotecc/cmd/tcpclient"
	"strings"
)

const DEBUG = false

type HWInterface struct {
	clients map[string]*tcpclient.Client
	PowerOn bool
	Online  bool
}

func NewHWInterface() *HWInterface {
	mpdAddr := storage.AppSettings.Server.IPAddr + ":" + storage.AppSettings.Server.MPDPort
	ctrlAddr := storage.AppSettings.Server.IPAddr + ":" + storage.AppSettings.Server.CtrlPort
	return &HWInterface{
		clients: map[string]*tcpclient.Client{
			"mpd":  tcpclient.New(tcpclient.Conf{Addr: mpdAddr, SingleRequest: false}),
			"ctrl": tcpclient.New(tcpclient.Conf{Addr: ctrlAddr, SingleRequest: true}),
		},
	}
}

func (hw *HWInterface) Request(server, command string) ([]string, error) {
	if DEBUG {
		switch command {
		case "status", "currentsong", "check_extpower":
		default:
			fmt.Printf("Rrq [%s]: [%s]\n", server, command)
		}
	}
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
