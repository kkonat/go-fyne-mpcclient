package storage

import (
	"fmt"

	"github.com/spf13/viper"
)

type ServerData struct {
	IPAddr     string
	MPDPort    string
	CtrlPort   string
	showHWCtrl bool
}
type SettingsData struct {
	Server ServerData
}

var AppSettings *SettingsData

func Init() {
	viper.SetConfigName("config.yml")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	AppSettings = &SettingsData{}

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Error reading config file, %s", err))
	}
	err := viper.Unmarshal(AppSettings)
	if err != nil {
		panic(fmt.Sprintf("Unable to decode into struct, %v", err))
	}
}

func Finalize() {
	err := viper.WriteConfig()
	if err != nil {
		panic(fmt.Sprintf("Unable to save config: %v", err))
	}
}
