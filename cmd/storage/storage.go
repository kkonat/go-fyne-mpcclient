package storage

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const CONFIGFILE = "remotecc-config.yml"

type ServerData struct {
	IPAddr   string
	MPDPort  string
	CtrlPort string
}
type SettingsData struct {
	Server     ServerData
	ShowHWCtrl bool
}

var AppSettings *SettingsData

func Init() {
	viper.SetConfigName(CONFIGFILE)
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
	// TODO : add message box if there's no CONFIGFILE
}

func SaveData() error {

	// due to some reasons, viper.WriteConfig() doesn't save the updated data structure,
	// instead, it saves viper config which was read into memory, i.e. from before the viper.Unmarshal()
	// so I had to manually Marshall the struct to yaml and save it,

	b, err := yaml.Marshal(AppSettings)
	if err != nil {
		return err
	}

	f, err := os.Create(CONFIGFILE)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(string(b))
	if err != nil {
		panic(fmt.Sprintf("Unable to save config: %v", err))
	}
	// log.Println(viper.ConfigFileUsed())

	return nil
}
