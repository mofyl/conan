package config

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"conan/utils"
)

const (
	CONFIGDIR        = "configfile"
	CONFIG_FILE_NAME = "config_test.yaml"
)

func init() {
	newFile("")
}

func newFile(base string) {
	if base == "" {
		//base = utils.GetProDirAbs() + CONFIGDIR
		base = GetConfigPath()
	}
	//
	viper.SetConfigName(CONFIG_FILE_NAME)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(base)
	loadFile()
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		switch in.Op {
		case fsnotify.Create, fsnotify.Write:
			loadFile()
		}
	})
}

func loadFile() {
	err := viper.ReadInConfig()
	if err != nil {
		panic("ReadConfig File " + err.Error())
	}
}

func SetFile(fileName string) {
	viper.SetConfigFile(fileName)
	loadFile()
}

func GetConfigPath() string {
	return utils.GetProDirAbs() + CONFIGDIR
}

func DecodeAll(fileName string, data interface{}) error {
	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")
	loadFile()

	return viper.Unmarshal(data)
}

func Decode(fileName string, key string, data interface{}) error {
	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")
	loadFile()

	if !viper.IsSet(key) {
		return errors.New("cur key not exist in any Config File")
	}
	return viper.UnmarshalKey(key, data)
}
